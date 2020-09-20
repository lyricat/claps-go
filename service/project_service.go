package service

import (
	"claps-test/dao"
	"claps-test/model"
	"claps-test/util"
	"github.com/gin-gonic/gin"
)

//通过projectName查询,查询某个项目的详情
func GetProjectByName(ctx *gin.Context, name string) (projectDetailInfo *map[string]interface{}, err *util.Err) {

	project, err1 := dao.GetProjectByName(name)
	if err1 != nil {
		err = util.NewErr(err1, util.ErrDataBase, "获取项目信息失败")
		return
	}

	repositoryDtos, err1 := dao.ListRepositoriesByProjectId(project.Id)
	if err1 != nil {
		err = util.NewErr(err1, util.ErrDataBase, "获取项目仓库失败")
		return
	}

	/*
		for i := range *repositoryDtos {
			repoInfo,err1 := GetRepositoryInfo(ctx,(*repositoryDtos)[i].Slug)
			if err1 != nil {
				err = util.NewErr(err1,util.ErrThirdParty,"获取项目仓库详细信息失败")
				return
			}
			(*repositoryDtos)[i].Forks = *repoInfo.ForksCount
			(*repositoryDtos)[i].Stars = *repoInfo.StargazersCount
			(*repositoryDtos)[i].Watchs = *repoInfo.WatchersCount
			(*repositoryDtos)[i].RepositoryUrl = *repoInfo.ArchiveURL
		}
	*/

	//mambers格式不同,删除project_id和userid字段
	members, err1 := dao.ListMembersByProjectName(project.Name)
	if err1 != nil {
		err = util.NewErr(err1, util.ErrDataBase, "获取项目成员失败")
		return
	}

	botDtos, err1 := dao.ListBotDtosByProjectId(project.Id)
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

//获取数据库中所有project
func ListProjectsAll() (projects *[]model.Project, err *util.Err) {
	projects, err1 := dao.ListProjectsAll()
	if err1 != nil {
		err = util.NewErr(err1, util.ErrDataBase, "获取所有项目失败")
	}
	return
}

//查询某用户的所有项目,获取数据库中所有project
func ListProjectsByUserId(userId int64) (projects *[]model.Project, err *util.Err) {
	projects, err1 := dao.ListProjectsByUserId(userId)
	if err1 != nil {
		err = util.NewErr(err1, util.ErrDataBase, "获取项目机器人失败")
	}

	return
}

func ListTransactionsByProjectName(name string) (transactions *[]model.Transaction, err *util.Err) {

	transactions, err1 := dao.ListTransactionsByProjectName(name)
	if err1 != nil {
		err = util.NewErr(err1, util.ErrDataBase, "获取项目获取捐赠记录失败")
	}
	return
}

func ListMembersByProjectName(projectName string) (members *[]model.User, err *util.Err) {
	members, err1 := dao.ListMembersByProjectName(projectName)
	if err1 != nil {
		err = util.NewErr(err1, util.ErrDataBase, "获取项目成员失败")
	}
	return
}
