package encoders

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	logging "github.com/remoteit/systemkit-logging"
	spec "github.com/remoteit/systemkit-service-spec"
)

var logTag = "UpStart-SERVICE"

// SERVICEToUpStart -
func SERVICEToUpStart(serviceSpec spec.SERVICE) (platformService string) {
	// for SystemD move everything into config.Executable
	if len(serviceSpec.Args) > 0 {
		serviceSpec.Executable = fmt.Sprintf(
			"%s %s",
			serviceSpec.Executable,
			strings.Join(serviceSpec.Args, " "),
		)
	}

	fileTemplate := template.Must(template.New("upstartFile").Parse(`# {{.Description}}

description     "{{.Name}}"

start on filesystem or runlevel [2345]
stop on runlevel [!2345]

#setuid username

# stop the respawn is process fails to start 5 times within 5 minutes
respawn
respawn limit 5 300
umask 022

console none

pre-start script
    test -x {{.Executable}} || { stop; exit 0; }
end script

# Start
exec {{.Executable}}
`))

	var buffer bytes.Buffer
	if err := fileTemplate.Execute(&buffer, serviceSpec); err != nil {
		logging.Errorf("%s: error generating file: %s", logTag, err.Error())
		return ""
	}

	return buffer.String()
}
