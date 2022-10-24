package flipper

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type gitShortLog struct {
	author     string
	authorDate time.Time
}

func GetGitShortLog(lineNo int, repoFolder string, filePath string) (*gitShortLog, error) {
	var out bytes.Buffer
	gitRangeFilter := fmt.Sprintf("%d,+1:%s", lineNo, filePath)
	cmd := exec.Command("git", "-C", repoFolder, "log", "--format=%aN <%ae>%n%aI", "--no-patch", "-L", gitRangeFilter)
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, err
	}
	// fmt.Printf("Execute: %s\n", strings.Join(cmd.Args, " "))

	infos := strings.Split(strings.ReplaceAll(out.String(), "\r\n", "\n"), "\n")
	if len(infos) < 2 {
		return nil, fmt.Errorf("Invalid git log result, len:%d", len(infos))
	}

	authorDate, err := time.Parse(time.RFC3339, infos[1])
	if err != nil {
		return nil, err
	}
	log := gitShortLog{
		author:     strings.TrimSpace(infos[0]),
		authorDate: authorDate,
	}

	return &log, nil
}
