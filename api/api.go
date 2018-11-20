package api

import (
	"fmt"
	communicator "git.costrategix.net/go/mavenlink-communicator/proto/mavenlink-communicator"
	"github.com/micro/go-log"
	"github.com/pkg/errors"
	"net/url"
	"strings"
)

// Fixed map
var endpoint = map[string]string{
	"workspaces":   "workspaces.json",
	"stories":      "stories.json",
	"time_entries": "time_entries.json",
	"users":        "users.json",
}

// MavenlinkApiInterface provides the interface definition for this service
type MavenlinkApiInterface interface {
	SetEnv(configuration *communicator.EnvironmentConfiguration) error
	GetProjects() ([]*communicator.Project, error)
	GetProject(keyOrId string) (*communicator.Project, error)
	GetTasksFromProjectId(keyOrId string) ([]*communicator.Task, error)
	GetSubTasksFromProjectId(workspace string, task string) ([]*communicator.Task, error)
	GetIssueTasksFromProjectId(keyOrId string, subTask string) ([]*communicator.Task, error)
	GetTimeEntriesFromProjectIdAndIssueTaskId(projectKeyOrId string, issueTaskKeyOrId string) ([]*communicator.Timeentry, error)
	GetUsersFromProjectId(projectKeyOrId string) ([]*communicator.User, error)
	GetUserFromProjectId(projectKeyOrId string, userId string) *communicator.User
	FormatErrors(err error, message string) *communicator.Error
}

// MavenlinkApi provides a concrete instance of the interface MavenlinkApiInterface
type MavenlinkApi struct {
	env *communicator.EnvironmentConfiguration
}

func (mavenlink *MavenlinkApi) SetEnv(configuration *communicator.EnvironmentConfiguration) error {
	if configuration == nil {
		return errors.New("No configurations detected")
	}
	mavenlink.env = configuration
	return nil
}

func (mavenlink *MavenlinkApi) FormatErrors(err error, message string) *communicator.Error {
	errResp := new(communicator.Error)
	errResp.Code = int32(400)
	if mavenlink.env.Debug == true {
		errResp.Description = err.Error()
	} else {
		errResp.Description = message
	}
	return errResp
}

// GetProjects is used to retrieve all the workspaces available in JIRA
func (mavenlink *MavenlinkApi) GetProjects() ([]*communicator.Project, error) {
	var workspacesResponse *communicator.MavenlinkWorkspacesResponse
	var projects []*communicator.Project
	var Url *url.URL
	Url, UrlErr := url.Parse(mavenlink.env.Url)
	if UrlErr != nil {
		return nil, errors.New("Failed to parse environment URL")
	}
	Url.Path += endpoint["workspaces"]
	token := mavenlink.env.Token
	apiErr := InsecureRequest(Url.String(), "GET", nil, token, &workspacesResponse)
	if apiErr != nil {
		if mavenlink.env.Debug == true {
			log.Logf("Error(API - %s) : %s\n", Url.String(), apiErr)
		}
		return nil, apiErr
	}
	if workspacesResponse == nil {
		var temp map[string]interface{}
		someErr := InsecureRequest(Url.String(), "GET", nil, token, &temp)
		if someErr != nil {
			return nil, errors.New("Failed to retrieve response from workspaces endpoint(Level 2)")
		}
		log.Logf("JSON response: %v\n", temp)
		return nil, errors.New("Failed to retrieve response from workspaces endpoint")
	}
	if workspacesResponse.Count > 0 {
		for key, workspace := range workspacesResponse.Workspaces {
			if fmt.Sprint(key) == workspace.Id {
				project := new(communicator.Project)
				project.Id = workspace.Id
				project.Title = workspace.Title
				project.Description = workspace.Description
				project.AccessLevel = workspace.AccessLevel
				project.AccountId = workspace.AccountId
				project.Archived = workspace.Archived
				project.Currency = workspace.Currency
				project.CurrencySymbol = workspace.CurrencySymbol
				project.CreatedAt = workspace.CreatedAt
				project.DueDate = workspace.DueDate
				project.EffectiveDueDate = workspace.EffectiveDueDate
				project.StartDate = workspace.StartDate
				project.UpdatedAt = workspace.UpdatedAt
				projects = append(projects, project)
			}
		}
		if len(projects) != int(workspacesResponse.Count) {
			return nil, errors.New(
				"Mismatch found between processed and retrieved count. Failed to retrieve all Projects!")
		}
	}
	return projects, nil
}

// GetProject is used to retrieve a single workspace from Mavenlink
func (mavenlink *MavenlinkApi) GetProject(keyOrId string) (*communicator.Project, error) {
	var workspacesResponse *communicator.MavenlinkWorkspacesResponse
	var project *communicator.Project
	var Url *url.URL
	Url, UrlErr := url.Parse(mavenlink.env.Url)
	if UrlErr != nil {
		return project, errors.New("Failed to parse environment URL")
	}
	Url.Path += endpoint["workspaces"]
	parameters := url.Values{}
	parameters.Add("only", fmt.Sprint(keyOrId))
	Url.RawQuery = parameters.Encode()
	token := mavenlink.env.Token
	apiErr := InsecureRequest(Url.String(), "GET", nil, token, &workspacesResponse)
	if apiErr != nil {
		if mavenlink.env.Debug == true {
			log.Logf("Error(API - %s) : %s\n", Url.String(), apiErr)
		}
		return project, apiErr
	}
	if workspacesResponse == nil || workspacesResponse.Workspaces == nil {
		return project, errors.New("Failed to retrieve response from workspaces endpoint")
	}
	for key, workspace := range workspacesResponse.Workspaces {
		if fmt.Sprint(key) == keyOrId {
			project := new(communicator.Project)
			project.Id = workspace.Id
			project.Title = workspace.Title
			project.Description = workspace.Description
			project.AccessLevel = workspace.AccessLevel
			project.AccountId = workspace.AccountId
			project.Archived = workspace.Archived
			project.Currency = workspace.Currency
			project.CurrencySymbol = workspace.CurrencySymbol
			project.CreatedAt = workspace.CreatedAt
			project.DueDate = workspace.DueDate
			project.EffectiveDueDate = workspace.EffectiveDueDate
			project.StartDate = workspace.StartDate
			project.UpdatedAt = workspace.UpdatedAt
			return project, nil
		}
	}
	return project, nil
}

// GetTasksFromProjectId is used to retrieve all the stories from a workspace in Mavenlink
func (mavenlink *MavenlinkApi) GetTasksFromProjectId(keyOrId string) ([]*communicator.Task, error) {
	var storiesResponse *communicator.MavenlinkStoriesResponse
	var tasks []*communicator.Task
	var Url *url.URL
	Url, UrlErr := url.Parse(mavenlink.env.Url)
	if UrlErr != nil {
		return tasks, errors.New("Failed to parse environment URL")
	}
	Url.Path += endpoint["stories"]
	parameters := url.Values{}
	parameters.Add("workspace_id", fmt.Sprint(keyOrId))
	parameters.Add("parents_only", "true")
	Url.RawQuery = parameters.Encode()
	token := mavenlink.env.Token
	apiErr := InsecureRequest(Url.String(), "GET", nil, token, &storiesResponse)
	if apiErr != nil {
		if mavenlink.env.Debug == true {
			log.Logf("Error(API - %s) : %s\n", Url.String(), apiErr)
		}
		return tasks, apiErr
	}
	if storiesResponse == nil || storiesResponse.Stories == nil {
		return tasks, errors.New("Failed to retrieve response from stories endpoint")
	}
	for _, story := range storiesResponse.Stories {
		if len(story.ParentId) < 1 {
			task := new(communicator.Task)
			task.Id = story.Id
			task.Title = story.Title
			task.Description = story.Description
			task.StoryType = story.StoryType
			task.Priority = story.Priority
			task.Archived = story.Archived
			task.WorkspaceId = story.WorkspaceId
			task.CreatorId = story.CreatorId
			task.ParentId = story.ParentId
			task.CreatedAt = story.CreatedAt
			task.DueDate = story.DueDate
			task.State = story.State
			task.StartDate = story.StartDate
			task.UpdatedAt = story.UpdatedAt
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

// GetSubTasksFromProjectId is used to retrieve all the stories from a workspace in Mavenlink
func (mavenlink *MavenlinkApi) GetSubTasksFromProjectId(workspace string, task string) ([]*communicator.Task, error) {
	var storiesResponse *communicator.MavenlinkStoriesResponse
	var tasks []*communicator.Task
	var Url *url.URL
	Url, UrlErr := url.Parse(mavenlink.env.Url)
	if UrlErr != nil {
		return tasks, errors.New("Failed to parse environment URL")
	}
	Url.Path += endpoint["stories"]
	parameters := url.Values{}
	parameters.Add("workspace_id", workspace)
	parameters.Add("with_parent_id", task)
	Url.RawQuery = parameters.Encode()
	token := mavenlink.env.Token
	apiErr := InsecureRequest(Url.String(), "GET", nil, token, &storiesResponse)
	if apiErr != nil {
		if mavenlink.env.Debug == true {
			log.Logf("Error(API - %s) : %s\n", Url.String(), apiErr)
		}
		return tasks, apiErr
	}
	if storiesResponse == nil || storiesResponse.Stories == nil {
		return tasks, errors.New("Failed to retrieve response from stories endpoint")
	}
	for _, story := range storiesResponse.Stories {
		if len(story.ParentId) > 0 && story.ParentId == task {
			task := new(communicator.Task)
			task.Id = story.Id
			task.Title = story.Title
			task.Description = story.Description
			task.StoryType = story.StoryType
			task.Priority = story.Priority
			task.Archived = story.Archived
			task.WorkspaceId = story.WorkspaceId
			task.CreatorId = story.CreatorId
			task.ParentId = story.ParentId
			task.CreatedAt = story.CreatedAt
			task.DueDate = story.DueDate
			task.State = story.State
			task.StartDate = story.StartDate
			task.UpdatedAt = story.UpdatedAt
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

// GetIssueTasksFromProjectId is used to retrieve all the stories from a workspace in Mavenlink
func (mavenlink *MavenlinkApi) GetIssueTasksFromProjectId(workspace string,
	subTask string) ([]*communicator.Task, error) {

	var storiesResponse *communicator.MavenlinkStoriesResponse
	var tasks []*communicator.Task
	var Url *url.URL
	Url, UrlErr := url.Parse(mavenlink.env.Url)
	if UrlErr != nil {
		return tasks, errors.New("Failed to parse environment URL")
	}
	Url.Path += endpoint["stories"]
	parameters := url.Values{}
	parameters.Add("workspace_id", workspace)
	parameters.Add("with_parent_id", subTask)
	parameters.Add("include", "assignees")
	Url.RawQuery = parameters.Encode()
	token := mavenlink.env.Token
	apiErr := InsecureRequest(Url.String(), "GET", nil, token, &storiesResponse)
	if apiErr != nil {
		if mavenlink.env.Debug == true {
			log.Logf("Error(API - %s) : %s\n", Url.String(), apiErr)
		}
		return tasks, apiErr
	}
	if storiesResponse == nil || storiesResponse.Stories == nil {
		return tasks, errors.New("Failed to retrieve response from stories endpoint")
	}
	for _, story := range storiesResponse.Stories {
		if len(story.ParentId) > 0 && story.ParentId == subTask {
			task := new(communicator.Task)
			task.Id = story.Id
			task.Title = story.Title
			task.Description = story.Description
			task.StoryType = story.StoryType
			task.Priority = story.Priority
			task.Archived = story.Archived
			task.WorkspaceId = story.WorkspaceId
			task.CreatorId = story.CreatorId
			task.ParentId = story.ParentId
			task.CreatedAt = story.CreatedAt
			task.DueDate = story.DueDate
			task.State = story.State
			task.StartDate = story.StartDate
			task.UpdatedAt = story.UpdatedAt
			if story.AssigneeIds != nil {
				for _, assignee := range story.AssigneeIds {
					task.User = mavenlink.GetUserFromProjectId(workspace, assignee)
				}
			}

			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

func (mavenlink *MavenlinkApi) GetTimeEntriesFromProjectIdAndIssueTaskId(projectKeyOrId string,
	issueTaskKeyOrId string) ([]*communicator.Timeentry, error) {

	var timeentriesResponse *communicator.MavenlinkTimeEntriesResponse
	var timeentries []*communicator.Timeentry
	var Url *url.URL
	Url, UrlErr := url.Parse(mavenlink.env.Url)
	if UrlErr != nil {
		return timeentries, errors.New("Failed to parse environment URL")
	}
	Url.Path += endpoint["time_entries"]
	parameters := url.Values{}
	parameters.Add("workspace_id", projectKeyOrId)
	Url.RawQuery = parameters.Encode()
	token := mavenlink.env.Token
	apiErr := InsecureRequest(Url.String(), "GET", nil, token, &timeentriesResponse)
	if apiErr != nil {
		if mavenlink.env.Debug == true {
			log.Logf("Error(API - %s) : %s\n", Url.String(), apiErr)
		}
		return timeentries, apiErr
	}
	if timeentriesResponse == nil || timeentriesResponse.TimeEntries == nil {
		return timeentries, errors.New("Failed to retrieve response from time entries endpoint")
	}
	for _, timeentry := range timeentriesResponse.TimeEntries {
		if strings.EqualFold(issueTaskKeyOrId, timeentry.StoryId) {
			timeentryWithUser := new(communicator.Timeentry)
			timeentryWithUser.Id = timeentry.Id
			timeentryWithUser.DatePerformed = timeentry.DatePerformed
			timeentryWithUser.TimeInMinutes = timeentry.TimeInMinutes
			timeentryWithUser.Notes = timeentry.Notes
			timeentryWithUser.WorkspaceId = timeentry.WorkspaceId
			timeentryWithUser.StoryId = timeentry.StoryId
			timeentryWithUser.CreatedAt = timeentry.CreatedAt
			timeentryWithUser.UpdatedAt = timeentry.UpdatedAt
			timeentryWithUser.User = mavenlink.GetUserFromProjectId(projectKeyOrId, timeentry.UserId)
			timeentries = append(timeentries, timeentryWithUser)
		}
	}
	return timeentries, nil
}

func (mavenlink *MavenlinkApi) GetUsersFromProjectId(projectKeyOrId string) ([]*communicator.User, error) {
	var usersResponse *communicator.MavenlinkUsersResponse
	var users []*communicator.User
	var Url *url.URL
	Url, UrlErr := url.Parse(mavenlink.env.Url)
	if UrlErr != nil {
		return users, errors.New("Failed to parse environment URL")
	}
	Url.Path += endpoint["users"]
	parameters := url.Values{}
	parameters.Add("participant_in", projectKeyOrId)
	Url.RawQuery = parameters.Encode()
	token := mavenlink.env.Token
	apiErr := InsecureRequest(Url.String(), "GET", nil, token, &usersResponse)
	if apiErr != nil {
		if mavenlink.env.Debug == true {
			log.Logf("Error(API - %s) : %s\n", Url.String(), apiErr)
		}
		return users, apiErr
	}
	if usersResponse == nil || usersResponse.Users == nil {
		return users, errors.New("Failed to retrieve response from users endpoint")
	}
	for _, user := range usersResponse.Users {
		formattedUser := new(communicator.User)
		formattedUser.Id = user.Id
		formattedUser.FullName = user.FullName
		formattedUser.EmailAddress = user.EmailAddress
		formattedUser.Headline = user.Headline
		formattedUser.AccountId = user.AccountId
		users = append(users, formattedUser)
	}
	return users, nil
}

func (mavenlink *MavenlinkApi) GetUserFromProjectId(projectKeyOrId string, userId string) *communicator.User {
	var theUser *communicator.User
	users, usersErr := mavenlink.GetUsersFromProjectId(projectKeyOrId)
	if usersErr != nil {
		return theUser
	}
	for _, user := range users {
		if user.Id == userId {
			formattedUser := new(communicator.User)
			formattedUser.Id = user.Id
			formattedUser.FullName = user.FullName
			formattedUser.EmailAddress = user.EmailAddress
			formattedUser.Headline = user.Headline
			formattedUser.AccountId = user.AccountId
			theUser = formattedUser
		}
	}
	return theUser
}
