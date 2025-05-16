package convert

import (
	db "go-cs/internal/bean/biz"
	domain "go-cs/internal/domain/space_work_version"
)

func SpaceWorkVersionEntityToPo(version *domain.SpaceWorkVersion) *db.SpaceWorkVersion {
	// 转换为数据库对象
	return &db.SpaceWorkVersion{
		Id:            version.Id,
		SpaceId:       version.SpaceId,
		VersionKey:    version.VersionKey,
		VersionName:   version.VersionName,
		VersionStatus: version.VersionStatus,
		Remark:        version.Remark,
		Ranking:       version.Ranking,
		CreatedAt:     version.CreatedAt,
		UpdatedAt:     version.UpdatedAt,
		DeletedAt:     version.DeletedAt,
	}
}

func SpaceWorkVersionEntityToPos(versions []*domain.SpaceWorkVersion) []*db.SpaceWorkVersion {
	var versionsPo []*db.SpaceWorkVersion
	for _, version := range versions {
		versionsPo = append(versionsPo, SpaceWorkVersionEntityToPo(version))
	}
	return versionsPo
}

func SpaceWorkVersionPoToEntity(version *db.SpaceWorkVersion) *domain.SpaceWorkVersion {
	return &domain.SpaceWorkVersion{
		Id:            version.Id,
		SpaceId:       version.SpaceId,
		VersionKey:    version.VersionKey,
		VersionName:   version.VersionName,
		VersionStatus: version.VersionStatus,
		Remark:        version.Remark,
		Ranking:       version.Ranking,
		CreatedAt:     version.CreatedAt,
		UpdatedAt:     version.UpdatedAt,
		DeletedAt:     version.DeletedAt,
	}
}

func SpaceWorkVersionPoToEntitys(versions []*db.SpaceWorkVersion) []*domain.SpaceWorkVersion {
	var versionsPo []*domain.SpaceWorkVersion
	for _, version := range versions {
		versionsPo = append(versionsPo, SpaceWorkVersionPoToEntity(version))
	}
	return versionsPo
}
