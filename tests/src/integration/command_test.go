// +build native

package tests

import (
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
	"github.com/openwhisk/openwhisk-cli/tests/src/integration/common"
)

var wsk *common.Wsk = common.NewWsk()
var tmpProp = os.Getenv("GOPATH") + "/src/github.com/openwhisk/openwhisk-cli/wskprops.tmp"
var invalidArgs []common.InvalidArg
var invalidParamMsg = "Arguments for '-p' must be a key/value pair"
var invalidAnnotMsg = "Arguments for '-a' must be a key/value pair"
var invalidParamFileMsg = "An argument must be provided for '-P'"
var invalidAnnotFileMsg = "An argument must be provided for '-A'"

var emptyFile = common.GetTestActionFilename("emtpy.js")
var helloFile = common.GetTestActionFilename("hello.js")
var missingFile = "notafile"
var emptyFileMsg = "File '" + emptyFile + "' is not a valid file or it does not exist"
var missingFileMsg = "File '" + missingFile + "' is not a valid file or it does not exist"

// Test case to check if the binary exits.
func TestWskExist(t *testing.T) {
	assert.True(t, wsk.Exists(), "The binary should exist.")
}

func TestHelpUsageInfoCommand(t *testing.T) {
	stdout, err := wsk.RunCommand("-h")
	assert.Equal(t, nil, err, "The command -h failed to run.")
	assert.Contains(t, string(stdout), "Usage:", "The output of the command -h does not contain \"Usage\".")
	assert.Contains(t, string(stdout), "Flags:", "The output of the command -h does not contain \"Flags\".")
	assert.Contains(t, string(stdout), "Available Commands:",
		"The output of the command -h does not contain \"Available Commands\".")
	assert.Contains(t, string(stdout), "--help", "The output of the command -h does not contain \"--help\".")
}

func TestHelpUsageInfoCommandLanguage(t *testing.T) {
	os.Setenv("LANG", "de_DE")
	assert.Equal(t, "de_DE", os.Getenv("LANG"), "The environment variable LANG has not been set to de_DE.")
	TestHelpUsageInfoCommand(t)
}

func TestShowCLIBuildVersion(t *testing.T) {
	stdout, err := wsk.RunCommand("property", "get", "--cliversion")
	assert.Equal(t, nil, err, "The command property get --cliversion failed to run.")
	assert.Contains(t, string(stdout), "whisk CLI version",
		"The output of the command property get --cliversion does not contain \"whisk CLI version\".")
}

func TestShowAPIVersion(t *testing.T) {
	stdout, err := wsk.RunCommand("property", "get", "--apiversion")
	assert.Equal(t, nil, err, "The command property get --apiversion failed to run.")
	assert.Contains(t, string(stdout), "whisk API version",
		"The output of the command property get --apiversion does not contain \"whisk API version\".")
}

// Test case to verify the default namespace _.
func TestDefaultNamespace(t *testing.T) {
	common.CreateFile(tmpProp)
	common.WriteFile(tmpProp, []string{"NAMESPACE="})

	os.Setenv("WSK_CONFIG_FILE", tmpProp)
	assert.Equal(t, os.Getenv("WSK_CONFIG_FILE"), tmpProp, "The environment variable WSK_CONFIG_FILE has not been set.")

	stdout, err := wsk.RunCommand("property", "get", "-i", "--namespace")
	assert.Equal(t, nil, err, "The command property get -i --namespace failed to run.")
	assert.Contains(t, common.RemoveRedundentSpaces(string(stdout)), "whisk namespace _",
		"The output of the command does not contain \"whisk namespace _\".")
	common.DeleteFile(tmpProp)
}

// Test case to validate default property values.
func TestValidateDefaultProperties(t *testing.T) {
	common.CreateFile(tmpProp)

	os.Setenv("WSK_CONFIG_FILE", tmpProp)
	assert.Equal(t, os.Getenv("WSK_CONFIG_FILE"), tmpProp, "The environment variable WSK_CONFIG_FILE has not been set.")

	stdout, err := wsk.RunCommand("property", "unset", "--auth", "--apihost", "--apiversion", "--namespace")
	assert.Equal(t, nil, err, "The command property unset failed to run.")
	outputString := string(stdout)
	assert.Contains(t, outputString, "ok: whisk auth unset",
		"The output of the command does not contain \"ok: whisk auth unset\".")
	assert.Contains(t, outputString, "ok: whisk API host unset",
		"The output of the command does not contain \"ok: whisk API host unset\".")
	assert.Contains(t, outputString, "ok: whisk API version unset",
		"The output of the command does not contain \"ok: whisk API version unset\".")
	assert.Contains(t, outputString, "ok: whisk namespace unset",
		"The output of the command does not contain \"ok: whisk namespace unset\".")

	stdout, err = wsk.RunCommand("property", "get", "--auth")
	assert.Equal(t, nil, err, "The command property get --auth failed to run.")
	assert.Equal(t, common.RemoveRedundentSpaces(string(stdout)), "whisk auth",
		"The output of the command does not equal to \"whisk auth\".")

	stdout, err = wsk.RunCommand("property", "get", "--apihost")
	assert.Equal(t, nil, err, "The command property get --apihost failed to run.")
	assert.Equal(t, common.RemoveRedundentSpaces(string(stdout)), "whisk API host",
		"The output of the command does not equal to \"whisk API host\".")

	stdout, err = wsk.RunCommand("property", "get", "--namespace")
	assert.Equal(t, nil, err, "The command property get --namespace failed to run.")
	assert.Equal(t, common.RemoveRedundentSpaces(string(stdout)), "whisk namespace _",
		"The output of the command does not equal to \"whisk namespace _\".")

	common.DeleteFile(tmpProp)
}

// Test case to set auth in property file.
func TestSetAuth(t *testing.T) {
	common.CreateFile(tmpProp)

	os.Setenv("WSK_CONFIG_FILE", tmpProp)
	assert.Equal(t, os.Getenv("WSK_CONFIG_FILE"), tmpProp, "The environment variable WSK_CONFIG_FILE has not been set.")

	_, err := wsk.RunCommand("property", "set", "--auth", "testKey")
	assert.Equal(t, nil, err, "The command property set --auth testKey failed to run.")
	output := common.ReadFile(tmpProp)
	assert.Contains(t, output, "AUTH=testKey",
		"The wsk property file does not contain \"AUTH=testKey\".")
	common.DeleteFile(tmpProp)
}

// Test case to set multiple property values with single command.
func TestSetMultipleValues(t *testing.T) {
	common.CreateFile(tmpProp)

	os.Setenv("WSK_CONFIG_FILE", tmpProp)
	assert.Equal(t, os.Getenv("WSK_CONFIG_FILE"), tmpProp, "The environment variable WSK_CONFIG_FILE has not been set.")

	_, err := wsk.RunCommand("property", "set", "--auth", "testKey", "--apihost", "openwhisk.ng.bluemix.net",
		"--apiversion", "v1")
	assert.Equal(t, nil, err, "The command property set --auth --apihost --apiversion failed to run.")
	output := common.ReadFile(tmpProp)
	assert.Contains(t, output, "AUTH=testKey", "The wsk property file does not contain \"AUTH=testKey\".")
	assert.Contains(t, output, "APIHOST=openwhisk.ng.bluemix.net",
		"The wsk property file does not contain \"APIHOST=openwhisk.ng.bluemix.net\".")
	assert.Contains(t, output, "APIVERSION=v1", "The wsk property file does not contain \"APIVERSION=v1\".")
	common.DeleteFile(tmpProp)
}

// Test case to reject bad command.
func TestRejectBadComm(t *testing.T) {
	common.CreateFile(tmpProp)

	os.Setenv("WSK_CONFIG_FILE", tmpProp)
	assert.Equal(t, os.Getenv("WSK_CONFIG_FILE"), tmpProp, "The environment variable WSK_CONFIG_FILE has not been set.")

	stdout, err := wsk.RunCommand("bogus")
	assert.NotEqual(t, nil, err, "The command bogus should fail to run.")
	assert.Contains(t, string(stdout), "Run 'wsk --help' for usage",
		"The output of the command does not contain \"Run 'wsk --help' for usage\".")
	common.DeleteFile(tmpProp)
}

// Test case to reject a command when the API host is not set.
func TestRejectCommAPIHostNotSet(t *testing.T) {
	common.CreateFile(tmpProp)

	os.Setenv("WSK_CONFIG_FILE", tmpProp)
	assert.Equal(t, os.Getenv("WSK_CONFIG_FILE"), tmpProp, "The environment variable WSK_CONFIG_FILE has not been set.")

	stdout, err := wsk.RunCommand("property", "get")
	assert.NotEqual(t, nil, err, "The command property get --apihost --apiversion should fail to run.")
	assert.Contains(t, common.RemoveRedundentSpaces(string(stdout)),
		"The API host is not valid: An API host must be provided",
		"The output of the command does not contain \"The API host is not valid: An API host must be provided\".")
	common.DeleteFile(tmpProp)
}

func initInvalidArgsNotEnoughParamsArgs() {
	invalidArgs = []common.InvalidArg{
		common.InvalidArg{
			Cmd: []string{"action", "create", "actionName", "-p"},
			Err: invalidParamMsg,
		},
		common.InvalidArg{
			Cmd: []string{"action", "create", "actionName", "-p", "key"},
			Err: invalidParamMsg,
		},
		common.InvalidArg{
			Cmd: []string{"action", "create", "actionName", "-P"},
			Err: invalidParamFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"action", "update", "actionName", "-p"},
			Err: invalidParamMsg,
		},
		common.InvalidArg{
			Cmd: []string{"action", "update", "actionName", "-p", "key"},
			Err: invalidParamMsg,
		},
		common.InvalidArg{
			Cmd: []string{"action", "update", "actionName", "-P"},
			Err: invalidParamFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"action", "invoke", "actionName", "-p"},
			Err: invalidParamMsg,
		},
		common.InvalidArg{
			Cmd: []string{"action", "invoke", "actionName", "-p", "key"},
			Err: invalidParamMsg,
		},
		common.InvalidArg{
			Cmd: []string{"action", "invoke", "actionName", "-P"},
			Err: invalidParamFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"action", "create", "actionName", "-a"},
			Err: invalidAnnotMsg,
		},
		common.InvalidArg{
			Cmd: []string{"action", "create", "actionName", "-a", "key"},
			Err: invalidAnnotMsg,
		},
		common.InvalidArg{
			Cmd: []string{"action", "create", "actionName", "-A"},
			Err: invalidAnnotFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"action", "update", "actionName", "-a"},
			Err: invalidAnnotMsg,
		},
		common.InvalidArg{
			Cmd: []string{"action", "update", "actionName", "-a", "key"},
			Err: invalidAnnotMsg,
		},
		common.InvalidArg{
			Cmd: []string{"action", "update", "actionName", "-A"},
			Err: invalidAnnotFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"action", "invoke", "actionName", "-a"},
			Err: invalidAnnotMsg,
		},
		common.InvalidArg{
			Cmd: []string{"action", "invoke", "actionName", "-a", "key"},
			Err: invalidAnnotMsg,
		},
		common.InvalidArg{
			Cmd: []string{"action", "invoke", "actionName", "-A"},
			Err: invalidAnnotFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"package", "create", "packageName", "-p"},
			Err: invalidParamMsg,
		},
		common.InvalidArg{
			Cmd: []string{"package", "create", "packageName", "-p", "key"},
			Err: invalidParamMsg,
		},
		common.InvalidArg{
			Cmd: []string{"package", "create", "packageName", "-P"},
			Err: invalidParamFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"package", "update", "packageName", "-p"},
			Err: invalidParamMsg,
		},
		common.InvalidArg{
			Cmd: []string{"package", "update", "packageName", "-p", "key"},
			Err: invalidParamMsg,
		},
		common.InvalidArg{
			Cmd: []string{"package", "update", "packageName", "-P"},
			Err: invalidParamFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"package", "bind", "packageName", "boundPackageName", "-p"},
			Err: invalidParamMsg,
		},
		common.InvalidArg{
			Cmd: []string{"package", "bind", "packageName", "boundPackageName", "-p", "key"},
			Err: invalidParamMsg,
		},
		common.InvalidArg{
			Cmd: []string{"package", "bind", "packageName", "boundPackageName", "-P"},
			Err: invalidParamFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"package", "create", "packageName", "-a"},
			Err: invalidAnnotMsg,
		},
		common.InvalidArg{
			Cmd: []string{"package", "create", "packageName", "-a", "key"},
			Err: invalidAnnotMsg,
		},
		common.InvalidArg{
			Cmd: []string{"package", "create", "packageName", "-A"},
			Err: invalidAnnotFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"package", "update", "packageName", "-a"},
			Err: invalidAnnotMsg,
		},
		common.InvalidArg{
			Cmd: []string{"package", "update", "packageName", "-a", "key"},
			Err: invalidAnnotMsg,
		},
		common.InvalidArg{
			Cmd: []string{"package", "update", "packageName", "-A"},
			Err: invalidAnnotFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"package", "bind", "packageName", "boundPackageName", "-a"},
			Err: invalidAnnotMsg,
		},
		common.InvalidArg{
			Cmd: []string{"package", "bind", "packageName", "boundPackageName", "-a", "key"},
			Err: invalidAnnotMsg,
		},
		common.InvalidArg{
			Cmd: []string{"package", "bind", "packageName", "boundPackageName", "-A"},
			Err: invalidAnnotFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"trigger", "create", "triggerName", "-p"},
			Err: invalidParamMsg,
		},
		common.InvalidArg{
			Cmd: []string{"trigger", "create", "triggerName", "-p", "key"},
			Err: invalidParamMsg,
		},
		common.InvalidArg{
			Cmd: []string{"trigger", "create", "triggerName", "-P"},
			Err: invalidParamFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"trigger", "update", "triggerName", "-p"},
			Err: invalidParamMsg,
		},
		common.InvalidArg{
			Cmd: []string{"trigger", "update", "triggerName", "-p", "key"},
			Err: invalidParamMsg,
		},
		common.InvalidArg{
			Cmd: []string{"trigger", "update", "triggerName", "-P"},
			Err: invalidParamFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"trigger", "fire", "triggerName", "-p"},
			Err: invalidParamMsg,
		},
		common.InvalidArg{
			Cmd: []string{"trigger", "fire", "triggerName", "-p", "key"},
			Err: invalidParamMsg,
		},
		common.InvalidArg{
			Cmd: []string{"trigger", "fire", "triggerName", "-P"},
			Err: invalidParamFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"trigger", "create", "triggerName", "-a"},
			Err: invalidAnnotMsg,
		},
		common.InvalidArg{
			Cmd: []string{"trigger", "create", "triggerName", "-a", "key"},
			Err: invalidAnnotMsg,
		},
		common.InvalidArg{
			Cmd: []string{"trigger", "create", "triggerName", "-A"},
			Err: invalidAnnotFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"trigger", "update", "triggerName", "-a"},
			Err: invalidAnnotMsg,
		},
		common.InvalidArg{
			Cmd: []string{"trigger", "update", "triggerName", "-a", "key"},
			Err: invalidAnnotMsg,
		},
		common.InvalidArg{
			Cmd: []string{"trigger", "update", "triggerName", "-A"},
			Err: invalidAnnotFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"trigger", "fire", "triggerName", "-a"},
			Err: invalidAnnotMsg,
		},
		common.InvalidArg{
			Cmd: []string{"trigger", "fire", "triggerName", "-a", "key"},
			Err: invalidAnnotMsg,
		},
		common.InvalidArg{
			Cmd: []string{"trigger", "fire", "triggerName", "-A"},
			Err: invalidAnnotFileMsg,
		},
	}
}

func initInvalidArgsMissingInvalidParamsAnno(){
	invalidArgs = []common.InvalidArg{
		common.InvalidArg{
			Cmd: []string{"action", "create", "actionName", helloFile, "-P", emptyFile},
			Err: emptyFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"action", "update", "actionName", helloFile, "-P", emptyFile},
			Err: emptyFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"action", "invoke", "actionName", "-P", emptyFile},
			Err: emptyFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"action", "create", "actionName", "-P", emptyFile},
			Err: emptyFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"action", "update", "actionName", "-P", emptyFile},
			Err: emptyFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"action", "invoke", "actionName", "-P", emptyFile},
			Err: emptyFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"package", "create", "packageName", "-P", emptyFile},
			Err: emptyFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"package", "update", "packageName", "-P", emptyFile},
			Err: emptyFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"package", "bind", "packageName", "boundPackageName", "-P", emptyFile},
			Err: emptyFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"package", "create", "packageName", "-P", emptyFile},
			Err: emptyFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"package", "update", "packageName", "-P", emptyFile},
			Err: emptyFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"package", "bind", "packageName", "boundPackageName", "-P", emptyFile},
			Err: emptyFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"trigger", "create", "triggerName", "-P", emptyFile},
			Err: emptyFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"trigger", "update", "triggerName", "-P", emptyFile},
			Err: emptyFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"trigger", "fire", "triggerName", "-P", emptyFile},
			Err: emptyFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"trigger", "create", "triggerName", "-P", emptyFile},
			Err: emptyFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"trigger", "update", "triggerName", "-P", emptyFile},
			Err: emptyFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"trigger", "fire", "triggerName", "-P", emptyFile},
			Err: emptyFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"action", "create", "actionName", helloFile, "-A", missingFile},
			Err: missingFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"action", "update", "actionName", helloFile, "-A", missingFile},
			Err: missingFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"action", "invoke", "actionName", "-A", missingFile},
			Err: missingFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"action", "create", "actionName", "-A", missingFile},
			Err: missingFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"action", "update", "actionName", "-A", missingFile},
			Err: missingFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"action", "invoke", "actionName", "-A", missingFile},
			Err: missingFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"package", "create", "packageName", "-A", missingFile},
			Err: missingFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"package", "update", "packageName", "-A", missingFile},
			Err: missingFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"package", "bind", "packageName", "boundPackageName", "-A", missingFile},
			Err: missingFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"trigger", "create", "triggerName", "-A", missingFile},
			Err: missingFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"trigger", "update", "triggerName", "-A", missingFile},
			Err: missingFileMsg,
		},
		common.InvalidArg{
			Cmd: []string{"trigger", "fire", "triggerName", "-A", missingFile},
			Err: missingFileMsg,
		},
	}
}

// Test case to reject commands that are executed with not enough param or annotation arguments.
func TestRejectCommandsNotEnoughParamsArgs(t *testing.T) {
	initInvalidArgsNotEnoughParamsArgs()
	for _, invalidArg := range invalidArgs {
		stdout, err := wsk.RunCommand(invalidArg.Cmd...)
		outputString := string(stdout)
		assert.NotEqual(t, nil, err, "The command should fail to run.")
		assert.Equal(t, "exit status 1", err.Error(), "The error should be exit status 1.")
		assert.Contains(t, outputString, invalidArg.Err,
			"The output of the command does not contain " + invalidArg.Err)
		assert.Contains(t, outputString, "Run 'wsk --help' for usage",
			"The output of the command does not contain \"Run 'wsk --help' for usage\".")
	}
}

// Test case to reject commands that are executed with a missing or invalid parameter or annotation file.
func TestRejectCommandsMissingIvalidParamsAnno(t *testing.T) {
	initInvalidArgsMissingInvalidParamsAnno()
	for _, invalidArg := range invalidArgs {
		stdout, err := wsk.RunCommand(invalidArg.Cmd...)
		outputString := string(stdout)
		assert.NotEqual(t, nil, err, "The command should fail to run.")
		assert.Equal(t, "exit status 2", err.Error(), "The error should be exit status 1.")
		assert.Contains(t, outputString, invalidArg.Err,
			"The output of the command does not contain " + invalidArg.Err)
		assert.Contains(t, outputString, "Run 'wsk --help' for usage",
			"The output of the command does not contain \"Run 'wsk --help' for usage\".")
	}
}