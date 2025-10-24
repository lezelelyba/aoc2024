package tests

// func TestServerStartup(t *testing.T) {
// 	var commonArgs = []string{"run", "main.go"}
// 	// cmd := exec.Command("go", "run", "../main.go")
//
// 	cases := []struct {
// 		name string
// 		args []string
// 	}{
// 		{"default", nil},
// 	}
//
// 	for _, c := range cases {
// 		t.Run(c.name, func(t *testing.T) {
// 			cmdArgs := append(commonArgs, c.args...)
// 			cmd := exec.Command("go", cmdArgs...)
// 			cmd.Dir = ".."
// 			out, err := cmd.CombinedOutput()
// 			if err != nil {
// 				t.Errorf("run failed %v\n%s", err, out)
// 			}
// 			t.Logf("output:\n%s", out)
// 			cmd.Process.Kill()
// 		})
// 	}
// }
