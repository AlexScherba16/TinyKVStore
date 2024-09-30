package application

// IApplication defines obligatory interface for application.
// core methods for starting and gracefully shutting down the application.
//
// Methods:
//
//   - Run() error: Starts the application.
//     Returns an error if something goes wrong during execution.
//
//   - Shutdown() error: Stops the application.
//     Returns an error if there are issues during the shutdown process.
type IApplication interface {
	Run() error
	Shutdown() error
}
