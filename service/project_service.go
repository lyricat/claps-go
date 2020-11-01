package service

import (
	"claps-test/model"
	"claps-test/util"
	log "github.com/sirupsen/logrus"
)

/**
 * @Description: 通过projectId查询,查询某个项目的详情
 * @param projectId
 * @return projectDetailInfo
 * @return err
 */
func GetProjectById(projectId int64) (projectDetailInfo *map[string]interface{}, err *util.Err) {

	project, err1 := model.PROJECT.GetProjectById(projectId)
	if err1 != nil {
		err = util.NewErr(err1, util.ErrDataBase, "获取项目信息失败")
		return
	}

	repositoryDtos, err1 := model.REPOSITORY.ListRepositoriesByProjectId(projectId)
	if err1 != nil {
		err = util.NewErr(err1, util.ErrDataBase, "获取项目仓库失败")
		return
	}

	//mambers格式不同,删除project_id和userid字段
	members, err1 := model.USER.ListMembersByProjectId(projectId)
	if err1 != nil {
		err = util.NewErr(err1, util.ErrDataBase, "获取项目成员失败")
		return
	}

	botDtos, err1 := model.BOTDTO.ListBotDtosByProjectId(projectId)
	if err1 != nil {
		err = util.NewErr(err1, util.ErrDataBase, "获取项目机器人失败")
		return
	}

	projectDetailInfo = &map[string]interface{}{
		"project":      project,
		"repositories": repositoryDtos,
		"members":      members,
		"botIds":       botDtos,
	}
	return
}

/**
 * @Description: 获取数据库中所有project,暂时弃用
 * @return projects
 * @return err
 */
func ListProjectsAll() (projects *[]model.Project, err *util.Err) {
	projects, err1 := model.PROJECT.ListProjectsAll()
	if err1 != nil {
		err = util.NewErr(err1, util.ErrDataBase, "获取所有项目失败")
	}
	return
}

/**
 * @Description: 通过query值获取project信息
 * @param q
 * @return projects
 * @return number
 * @return err
 */
func ListProjectsByQuery(q *model.PaginationQ) (projects *[]model.Project, number int, err *util.Err) {
	projects, number, err1 := model.PROJECT.ListProjectsByQuery(q)
	if err1 != nil {
		err = util.NewErr(err1, util.ErrDataBase, "根据query获取项目失败失败")
	}
	return
}

/**
 * @Description: 查询某用户的所有项目,获取数据库中所有project
 * @param userId
 * @return projects
 * @return err
 */
func ListProjectsByUserId(userId int64) (projects *[]model.Project, err *util.Err) {
	projects, err1 := model.PROJECT.ListProjectsByUserId(userId)
	if err1 != nil {
		err = util.NewErr(err1, util.ErrDataBase, "获取项目机器人失败")
	}

	return
}

/**
 * @Description: 获取数据库中对应projectId的所有transaction,暂时弃用
 * @param projectId
 * @return transactions
 * @return err
 */
func ListTransactionsByProjectId(projectId int64) (transactions *[]model.Transaction, err *util.Err) {

	transactions, err1 := model.TRANSACTION.ListTransactionsByProjectId(projectId)
	if err1 != nil {
		err = util.NewErr(err1, util.ErrDataBase, "获取项目获取捐赠记录失败")
	}
	return
}

/**
 * @Description: 通过projectId和query值获取transaction信息
 * @param projectId
 * @param q
 * @return transactions
 * @return number
 * @return err
 */
func ListTransactionsByProjectIdAndQuery(projectId int64, q *model.PaginationQ) (transactions *[]model.Transaction, number int, err *util.Err) {

	transactions, number, err1 := model.TRANSACTION.ListTransactionsByProjectIdAndQuery(projectId, q)
	if err1 != nil {
		err = util.NewErr(err1, util.ErrDataBase, "获取项目获取捐赠记录失败")
	}
	return
}

/**
 * @Description: 通过projectId获取对应项目成员的信息
 * @param projectId
 * @return members
 * @return err
 */
func ListMembersByProjectId(projectId int64) (members *[]model.User, err *util.Err) {
	members, err1 := model.USER.ListMembersByProjectId(projectId)
	if err1 != nil {
		err = util.NewErr(err1, util.ErrDataBase, "获取项目成员失败")
	}
	return
}

/**
 * @Description: TODO 获取项目badge,未完成
 * @param badge
 * @return err
 */
func GetProjectBadge(badge *model.Badge) (err *util.Err) {
	//compact
	//full

	fiat, err1 := model.FIAT.GetFiatByCode(badge.Code)
	if err1 != nil {
		err = util.NewErr(err1, util.ErrDataBase, "获取fiat失败")
	}
	log.Debug(fiat)
	return
}

/**
 * @Description:
 * @param projectId 项目Id
 * @return groupId 根据project_to_merico_group查到
 * @return err
 */
func GetGroupIdByProjectId(projectId int64) (groupId string, err *util.Err) {
	ptm, err1 := model.PROJECTTOMERICOGROUP.GetByProjectId(projectId)
	if err1 != nil {
		err = util.NewErr(err1, util.ErrDataBase, "DB get groupId error")
		return
	}
	groupId = ptm.MericoGroupId
	return
}
