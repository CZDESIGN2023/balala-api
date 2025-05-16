package biz

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	vo "go-cs/internal/bean/vo"
	"go-cs/internal/pkg/trans"
	"go-cs/internal/server/file/ffmpeg"
	"go-cs/internal/utils"
	"go-cs/internal/utils/rand"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	file_domain "go-cs/internal/domain/file_info"
	file_repo "go-cs/internal/domain/file_info/repo"
	file_service "go-cs/internal/domain/file_info/service"

	"github.com/go-kratos/kratos/v2/log"
)

type UploadUsecase struct {
	log             *log.Helper
	tm              trans.Transaction
	fileInfoRepo    file_repo.FileInfoRepo
	fileInfoService *file_service.FileInfoService
}

// NewUploadUsecase 初始化 UploadUsecase
func NewUploadUsecase(
	tm trans.Transaction,
	logger log.Logger,

	fileInfoRepo file_repo.FileInfoRepo,
	fileInfoService *file_service.FileInfoService,
) *UploadUsecase {
	moduleName := "UploadUsecase"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &UploadUsecase{
		log:             hlog,
		tm:              tm,
		fileInfoRepo:    fileInfoRepo,
		fileInfoService: fileInfoService,
	}
}

func (uc *UploadUsecase) GetLocalUserAvatarUploadDir(ctx context.Context, uid int64) (string, string) {
	nowYmd := time.Now().Format("20060102")

	filePath := path.Join("/", "avatar", nowYmd, strconv.FormatInt(uid, 10))
	uriPath := path.Join("/", "avatar", nowYmd, strconv.FormatInt(uid, 10))
	return filePath, uriPath
}

func (uc *UploadUsecase) GetLocalSpaceFileUploadDir(ctx context.Context, spaceId int64) (string, string) {
	nowYmd := time.Now().Format("20060102")

	filePath := path.Join("/", "space_file", nowYmd, strconv.FormatInt(spaceId, 10))
	uriPath := path.Join("/", "space_file", nowYmd, strconv.FormatInt(spaceId, 10))
	return filePath, uriPath
}

func (uc *UploadUsecase) GetLocalDefaultFileUploadDir(ctx context.Context, userId int64) (string, string) {
	nowYmd := time.Now().Format("20060102")

	filePath := path.Join("/", "user_default", nowYmd, strconv.FormatInt(userId, 10))
	uriPath := path.Join("/", "user_default", nowYmd, strconv.FormatInt(userId, 10))
	return filePath, uriPath
}

func (uc *UploadUsecase) DownloadAvatarImgToLocal(ctx context.Context, uid int64, avatarUrl string, downloadRootPath string) (string, error) {

	if avatarUrl == "" {
		return "", errors.New("头像地址为空")
	}

	ext := filepath.Ext(avatarUrl)

	filePath, fileUri := uc.GetLocalUserAvatarUploadDir(ctx, uid)
	err := os.MkdirAll(downloadRootPath+filePath, os.ModePerm)
	if err != nil {
		fmtErrMsg := fmt.Sprintf("创建文件夹失败: %v", err.Error())
		return "", errors.New(fmtErrMsg)
	}

	resp, err := http.Get(avatarUrl)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	// 创建一个文件用于保存
	fileName := "avatar_" + rand.S(5) + ext
	fullFileName := downloadRootPath + filePath + "/" + fileName
	out, err := os.Create(fullFileName)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	//保存用户头像Uri内容
	return fileUri + "/" + fileName, nil
}

// SaveFileToLocal 上传文件至本地
func (uc *UploadUsecase) SaveFileToLocal(ctx context.Context, in *vo.UploadFileToLocalVo) (*file_domain.FileInfo, error) {

	fileInfo := in.FileInfo
	dirPath := path.Join(in.FileRootPath, fileInfo.FileUploadPath)

	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return nil, err
	}

	// 获取文件扩展名
	ext := path.Ext(in.FileInfo.FileName)

	srcFile := in.FileInfo.File

	// 创建目标文件
	dstFile, err := os.CreateTemp(dirPath, "*")
	if err != nil {
		return nil, err
	}
	defer dstFile.Close()

	hash := md5.New()

	// 只需读取一次文件，同时进行 内容拷贝到目标文件 和 计算md5
	multiWriter := io.MultiWriter(hash, dstFile)
	if _, err := io.Copy(multiWriter, srcFile); err != nil {
		return nil, err
	}

	// 获取文件大小
	stat, _ := dstFile.Stat()
	fileSize := stat.Size()

	// 重命名目标文件
	md5Sum := hex.EncodeToString(hash.Sum(nil))
	fileName := fmt.Sprintf("%s_%d%s", md5Sum, time.Now().UnixMilli(), ext)
	localFilePath := filepath.Join(dirPath, fileName)
	if err := os.Rename(dstFile.Name(), localFilePath); err != nil {
		return nil, err
	}

	var fileUri = fileInfo.FileUri + "/" + fileName
	var coverUri = fileUri + ".jpg"

	// 生成封面路径
	if !ffmpeg.ExtractFirstFrame(localFilePath, localFilePath+".jpg") {
		coverUri = ""
	}

	// 获取元数据
	metaMap := GetFileMeta(localFilePath)

	//入库
	req := &file_service.CreateFileInfoRequest{
		Hash:         md5Sum,
		Name:         fileInfo.FileName,
		Uri:          fileUri,
		Cover:        coverUri,
		Typ:          0,
		Size:         fileSize,
		Owner:        in.OwnerId,
		Meta:         metaMap,
		UploadDomain: "local",
		UploadMd5:    md5Sum,
		UploadPath:   fileInfo.FileUploadPath,
	}

	res, err := uc.fileInfoService.CreateFileInfo(ctx, req)
	if err != nil {
		return nil, errors.New("file info save failed")
	}

	err = uc.fileInfoRepo.CreateFileInfo(ctx, res)
	if err != nil {
		return nil, errors.New("file info save failed")
	}

	return res, nil
}

func GetFileMeta(filePath string) map[string]string {
	var metaMap map[string]string
	meta, _ := ffmpeg.Probe(filePath)
	if meta != nil && len(meta.Streams) > 0 {
		metaMap = make(map[string]string)
		stream := meta.Streams[0]
		metaMap["codec_type"] = stream.CodecType
		duration := meta.Format.DurationSeconds

		switch stream.CodecType {
		case "video":
			metaMap["width"] = strconv.Itoa(stream.Width)
			metaMap["height"] = strconv.Itoa(stream.Height)
			metaMap["duration"] = cast.ToString(duration)
		case "audio":
			metaMap["duration"] = cast.ToString(duration)
		}
	}

	return metaMap
}
