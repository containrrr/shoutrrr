package cmd

const (
	//ExSuccess is the exit code that signals that everything went as expected
	ExSuccess = 0
	//ExUsage is the exit code that signals that the application was not started with the correct arguments
	ExUsage = 64
	//ExUnavailable is the exit code that signals that the application failed to perform the intended task
	ExUnavailable = 69
	//ExConfig is the exit code that signals that the task failed due to a configuration error
	ExConfig = 78
)
