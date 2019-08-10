workflow "Test" {
  on = "push"
  resolves = [
    "test",
    "lint",
    "notify"
  ]
}

action "build_ssh_server" {
  uses = "actions/docker/cli@master"
  args = "build -t ssh-server -f test/Dockerfile ."
}

action "run_ssh_server" {
  needs = ["build_ssh_server"]
  uses = "actions/docker/cli@master"
  args = "run -d --rm --name ssh-server -p 22:22 ssh-server"
}

action "test" {
  needs = ["run_ssh_server"]
  uses = "cedrickring/golang-action@1.2.0"
}

action "lint" {
  uses = "actions-contrib/golangci-lint@master"
  args = "run"
}

action "notify" {
  needs = [ "test" ]
  uses = "metrue/noticeme-github-action@master"
  secrets = ["NOTICE_ME_TOKEN"]
  args = ["Code Pushed in go-ssh-client repo"]
}
