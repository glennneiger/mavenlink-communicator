package main

import (
	API "github.com/desertjinn/mavenlink-communicator/api"
	LOG "github.com/desertjinn/mavenlink-communicator/log"
	communicator "github.com/desertjinn/mavenlink-communicator/proto/mavenlink-communicator"
	"github.com/kelseyhightower/envconfig"
	"github.com/micro/go-micro"
	//k8s "github.com/micro/kubernetes/go/micro"
	"golang.org/x/net/context"
	"log"
)

// Define the interface available in this service
type service struct {
	mavenlink API.MavenlinkApiInterface
}

// GetAllProjects can be used to retrieve the list of all available projects
func (s *service) GetAllProjects(ctx context.Context, req *communicator.Request, res *communicator.Response) error {
	// Retrieve all projects
	projects, err := s.mavenlink.GetProjects()
	if err != nil {
		return err
	}
	// Assign retrieved project data to response
	res.Projects = projects
	return nil
}

// GetProjectById can be used to retrieve a single project by ID from Mavenlink
func (s *service) GetProjectById(ctx context.Context, req *communicator.Request, res *communicator.Response) error {
	// Retrieve all projects
	project, err := s.mavenlink.GetProject(req.Workspace)
	if err != nil {
		return err
	}
	// Assign retrieved project data to response
	res.Project = project
	return nil
}

// GetTasksByProjectId can be used to retrieve all stories by workspace ID from Mavenlink
func (s *service) GetTasksByProjectId(ctx context.Context, req *communicator.Request, res *communicator.Response) error {
	// Retrieve all projects
	tasks, err := s.mavenlink.GetTasksFromProjectId(req.Workspace)
	if err != nil {
		return err
	}
	// Assign retrieved tasks to response
	res.Tasks = tasks
	return nil
}

// GetSubTasksByProjectId can be used to retrieve all stories by workspace ID from Mavenlink
func (s *service) GetSubTasksByParentTaskAndProjectId(ctx context.Context, req *communicator.Request, res *communicator.Response) error {
	// Retrieve all projects
	tasks, err := s.mavenlink.GetSubTasksFromProjectId(req.Workspace, req.Task)
	if err != nil {
		return err
	}
	// Assign retrieved tasks to response
	res.Tasks = tasks
	return nil
}

// GetTasksFromParentTaskAndProjectId can be used to retrieve all stories by workspace ID from Mavenlink
func (s *service) GetTasksBySubTaskParentTaskAndProjectId(ctx context.Context, req *communicator.Request, res *communicator.Response) error {
	// Retrieve all projects
	tasks, err := s.mavenlink.GetIssueTasksFromProjectId(req.Workspace, req.SubTask)
	if err != nil {
		return err
	}
	// Assign retrieved tasks to response
	res.Tasks = tasks
	return nil
}

// GetTasksFromParentTaskAndProjectId can be used to retrieve all stories by workspace ID from Mavenlink
func (s *service) GetTimeentries(ctx context.Context, req *communicator.Request, res *communicator.Response) error {
	// Retrieve all projects
	timeentries, err := s.mavenlink.GetTimeEntriesFromProjectIdAndIssueTaskId(req.Workspace, req.Task)
	if err != nil {
		return err
	}
	// Assign retrieved tasks to response
	res.Timeentries = timeentries
	return nil
}

// GetTasksFromParentTaskAndProjectId can be used to retrieve all stories by workspace ID from Mavenlink
func (s *service) GetUsers(ctx context.Context, req *communicator.Request, res *communicator.Response) error {
	// Retrieve all projects
	users, err := s.mavenlink.GetUsersFromProjectId(req.Workspace)
	if err != nil {
		return err
	}
	// Assign retrieved tasks to response
	res.Users = users
	return nil
}

// GetTasksFromParentTaskAndProjectId can be used to retrieve all stories by workspace ID from Mavenlink
func (s *service) GetUser(ctx context.Context, req *communicator.Request, res *communicator.Response) error {
	// Retrieve all projects
	user := s.mavenlink.GetUserFromProjectId(req.Workspace, req.KeyOrId)
	// Assign retrieved tasks to response
	res.User = user
	return nil
}

func main() {
	var env communicator.EnvironmentConfiguration
	// Retrieve environment configuration
	envConfigErr := envconfig.Process("mavenlink-communicator", &env)
	if envConfigErr != nil {
		log.Fatal(envConfigErr)
	}

	// Create an instance of the interface provided in this service
	// note: we're not setting env during initialization as the struct
	//       members are private
	mavenlink := &API.MavenlinkApi{}
	mavenlink.SetEnv(&env)

	// Create a new service
	srv := micro.NewService(
		// This name must match the package name given in the protobuf definition
		micro.Name("costrategix.service.mavenlink.communicator"),
		micro.Version("v1"),
		// Specify a log wrapper to log requests to this service in the console
		micro.WrapHandler(LOG.ConsoleLogWrapper),
	)
	// Init will parse the command line flags.
	srv.Init()

	// Register handler
	communicator.RegisterMavenlinkCommunicatorHandler(srv.Server(), &service{mavenlink})

	// Run the server
	if serverError := srv.Run(); serverError != nil {
		if env.Debug == true {
			log.Fatal(serverError)
		} else {
			log.Printf("Error initializing the service : %s", serverError.Error())
		}

	}
}
