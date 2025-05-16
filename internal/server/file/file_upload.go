package file

import (
	"context"
	"errors"
	"go-cs/internal/bean/vo"
	"io"
	"os"
	"path"
	"path/filepath"
)

func (s *Server) verify(ctx context.Context, uid int64, args VerifyArgs) (any, error) {
	// 检查文件是否存在
	var fileDirPath = path.Join(s.confFile.LocalPath, ".temp", args.FileHash)
	var finalFilePath = path.Join(s.confFile.LocalPath, ".merged", args.FileHash)

	if Exist(finalFilePath) {
		data, err := s.rapidUpload(ctx, uid, &Args{
			SpaceId:  args.SpaceId,
			Scene:    args.Scene,
			SubScene: args.SubScene,
			FileName: args.FileName,
			FileHash: args.FileHash,
		})
		if err != nil {
			return nil, err
		}

		return map[string]any{
			"shouldUpload": false,
			"uploadedList": []string{},
			"fileInfo":     data,
		}, nil
	}

	names := []string{}
	if Exist(fileDirPath) {
		names, _ = AllFileNamesOfDir(fileDirPath)
	}

	return map[string]any{
		"shouldUpload": true,
		"uploadedList": names,
	}, nil
}

func (s *Server) rapidChunk(ctx context.Context, uid int64, args Upload2Args) (any, error) {
	var fileDirPath = filepath.Join(s.confFile.LocalPath, ".temp", args.FileHash)

	err := os.MkdirAll(fileDirPath, os.ModePerm)
	if err != nil {
		return nil, err
	}

	dst, err := os.Create(path.Join(fileDirPath, args.ChunkHash))
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	src, err := args.ChunkFile.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func MergeFile(dirPath, finalFilePath string) (err error) {
	defer os.RemoveAll(dirPath)

	open, err := os.Open(dirPath)
	if err != nil {
		return err
	}
	defer open.Close()

	chunkPaths, err := AllFilePathOfDir(dirPath)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Dir(finalFilePath), 0755)
	if err != nil {
		return err
	}
	dst, err := os.Create(finalFilePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	SortPath(chunkPaths)
	for _, src := range chunkPaths {
		src, err := os.Open(src)
		if err != nil {
			return err
		}
		_, err = io.Copy(dst, src)
		if err != nil {
			src.Close()
			return err
		}

		src.Close()
	}

	return nil
}

func (s *Server) uploadHandler(ctx context.Context, uid int64, args UploadArgs) (any, error) {

	var fileSizeLimit int64 = MaxFileSize

	filePath := ""
	uriPath := ""

	//不同场景下的权限校验
	switch args.Scene {
	case "avatar":
		filePath, uriPath = s.uploadUsecase.GetLocalUserAvatarUploadDir(ctx, uid)

	case "space_file":
		member, err := s.spaceMemberRepo.GetSpaceMember(ctx, args.SpaceId, uid)
		if member == nil || err != nil { //成员不存在 不允许操作
			return nil, errors.New("无权限")
		}

		if args.SubScene == "attach" {
			fileSizeLimit = s.configRepo.AttachSize(ctx)
		}

		filePath, uriPath = s.uploadUsecase.GetLocalSpaceFileUploadDir(ctx, args.SpaceId)
	case "default":
		filePath, uriPath = s.uploadUsecase.GetLocalDefaultFileUploadDir(ctx, uid)
	default:
		return nil, errors.New("不支持的场景")
	}

	result := make(map[string]interface{})
	for _, file := range args.Files {
		// 检查文件大小
		if file.Size > fileSizeLimit {
			return nil, errors.New("文件大小超出限制")
		}
	}

	for _, file := range args.Files {
		open, _ := file.Open()

		uploadVo := &vo.UploadFileToLocalVo{
			FileRootPath: s.confFile.LocalPath,
			FileInfo: &vo.UploadFileInfoVo{
				File:           open,
				FileUri:        uriPath,
				FileUploadPath: filePath,
				FileName:       file.Filename,
			},
			OwnerId: uid,
		}

		var fileInfo map[string]any

		res, err := s.uploadUsecase.SaveFileToLocal(ctx, uploadVo)
		if err != nil {
			fileInfo = map[string]any{
				"status": "failed",
				"reason": err.Error(),
			}
		} else {
			fileInfo = map[string]any{
				"id":     res.Id,
				"hash":   res.UploadMd5,
				"size":   res.Size,
				"name":   res.Name,
				"o_name": file.Filename,
				"uri":    res.Uri,
				"cover":  res.Cover,
				"meta":   res.Meta,
				"status": "success",
			}
		}

		result[file.Filename] = fileInfo
	}

	return result, nil
}

func (s *Server) getPathByScene(ctx context.Context, scene, subScene string, spaceId int64) (filePath string, uriPath string, fileSizeLimit int64, err error) {
	uid := GetUserIdFromCtx(ctx)

	fileSizeLimit = MaxFileSize

	switch scene {
	case "avatar":
		filePath, uriPath = s.uploadUsecase.GetLocalUserAvatarUploadDir(ctx, uid)
	case "space_file":
		member, err := s.spaceMemberRepo.GetSpaceMember(ctx, spaceId, uid)
		if member == nil || err != nil { //成员不存在 不允许操作
			return "", "", 0, errors.New("无权限")
		}

		if subScene == "attach" {
			fileSizeLimit = s.configRepo.AttachSize(ctx)
		}

		filePath, uriPath = s.uploadUsecase.GetLocalSpaceFileUploadDir(ctx, spaceId)
	case "default":
		filePath, uriPath = s.uploadUsecase.GetLocalDefaultFileUploadDir(ctx, uid)
	default:
		return "", "", 0, errors.New("不支持的场景")
	}

	return
}

type Args struct {
	SpaceId  int64  `form:"spaceId"`
	Scene    string `form:"scene"`
	SubScene string `form:"subScene"`
	FileName string `form:"fileName"`
	FileHash string `form:"fileHash"`
}

func (s *Server) merge(ctx context.Context, uid int64, args MergeArgs) (data any, err error) {

	//不同场景下的权限校验
	filePath, uriPath, _, err := s.getPathByScene(ctx, args.Scene, args.SubScene, args.SpaceId)
	if err != nil {
		return nil, err
	}

	var fileDirPath = path.Join(s.confFile.LocalPath, ".temp", args.FileHash)
	var finalFilePath = path.Join(s.confFile.LocalPath, ".merged", args.FileHash)

	err = MergeFile(fileDirPath, finalFilePath)
	if err != nil {
		return nil, err
	}
	finalFile, _ := os.Open(finalFilePath)

	result := make(map[string]interface{})

	uploadVo := &vo.UploadFileToLocalVo{
		OwnerId:      uid,
		FileRootPath: s.confFile.LocalPath,
		FileInfo: &vo.UploadFileInfoVo{
			File:           finalFile,
			FileName:       args.FileName,
			FileUri:        uriPath,
			FileUploadPath: filePath,
		},
	}

	res, err := s.uploadUsecase.SaveFileToLocal(ctx, uploadVo)
	if err != nil {
		return nil, err
	}

	fileInfo := map[string]any{
		"id":          res.Id,
		"hash":        res.UploadMd5,
		"size":        res.Size,
		"name":        res.Name,
		"o_name":      args.FileName,
		"uri":         res.Uri,
		"cover":       res.Cover,
		"status":      "success",
		"meta":        res.Meta,
		"create_time": res.CreatedAt,
	}

	result[args.FileName] = fileInfo

	return result, nil
}

func (s *Server) rapidUpload(ctx context.Context, uid int64, args *Args) (data any, err error) {

	//不同场景下的权限校验
	filePath, uriPath, _, err := s.getPathByScene(ctx, args.Scene, args.SubScene, args.SpaceId)
	if err != nil {
		return nil, err
	}

	open, err := os.Open(path.Join(s.confFile.LocalPath, ".merged", args.FileHash))
	if err != nil {
		return nil, err
	}
	defer open.Close()

	result := make(map[string]interface{})

	uploadVo := &vo.UploadFileToLocalVo{
		OwnerId:      uid,
		FileRootPath: s.confFile.LocalPath,
		FileInfo: &vo.UploadFileInfoVo{
			File:           open,
			FileName:       args.FileName,
			FileUri:        uriPath,
			FileUploadPath: filePath,
		},
	}

	res, err := s.uploadUsecase.SaveFileToLocal(ctx, uploadVo)
	if err != nil {
		return nil, err
	}

	fileInfo := map[string]any{
		"id":          res.Id,
		"hash":        res.UploadMd5,
		"size":        res.Size,
		"name":        res.Name,
		"o_name":      args.FileName,
		"uri":         res.Uri,
		"cover":       res.Cover,
		"meta":        res.Meta,
		"status":      "success",
		"create_time": res.CreatedAt,
	}

	result[args.FileName] = fileInfo

	return result, nil
}
