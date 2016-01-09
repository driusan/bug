	commit string
	log    string
func (c HgCommit) CommitID() string {
	return c.commit
}
func (c HgCommit) LogMsg() string {
	return c.log
}
	return runCmd("hg", "log", "-p", "-g", "-r", c.commit, "--template={changelog}")

type HgTester struct {
	handler SCMHandler
	workdir string
}

func (t HgTester) GetLogs() ([]Commit, error) {
	logs, err := runCmd("hg", "log", "-r", ":", "--template", "{node} {desc}\\n")
	commits := make([]Commit, len(logMsgs)-1)
func (c HgTester) AssertStagingIndex(t *testing.T, f []FileStatus) {
    for _, file := range f {
        out, err := runCmd("hg", "status", file.Filename)
        if err != nil {
            t.Error("Could not get status of " + file.Filename)
        }

        // hg status doesn't include the working directory status
        expected := file.IndexStatus + " " + file.Filename + "\n"
        if out != expected {
            t.Error("Unexpected status. Got " + out + " not " + expected)
        }

    }
}

func (c HgTester) StageFile(file string) error {
	_, err := runCmd("hg", "add", file)
	return err
}
func (c *HgTester) Setup() error {
		c.workdir = dir
		os.Chdir(c.workdir)
	c.handler = HgManager{}

func (c HgTester) TearDown() {
	os.RemoveAll(c.workdir)
func (c HgTester) GetWorkDir() string {
	return c.workdir
}
func (c HgTester) AssertCleanTree(t *testing.T) {
		fmt.Printf("\"%s\"\n", out)
func (m HgTester) GetManager() SCMHandler {
	return m.handler
}
func TestHgBugRenameCommits(t *testing.T) {
	tester := HgTester{}
	expectedDiffs := []string{
		`diff --git a/issues/Test-bug/Description b/issues/Test-bug/Description
`, `diff --git a/issues/Renamed-bug/Description b/issues/Renamed-bug/Description