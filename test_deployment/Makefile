ARTIFACTS = \
  .terraform \
  terraform.tfstate \
  terraform.tfstate.backup \
  provision.retry \
  _init.ok \
  _server.ok \
  _provision.ok \
  _known_hosts \
  _ssh_config \
  _hosts \
  rtun \
  rtun.yml \
  rtun-server \
  rtun-server.yml

TF_OPTIONS = \
  -var-file .vars.json


.PHONY: all clean destroy mosh serve

all: _provision.ok
	@:

clean:
	rm -rf $(ARTIFACTS)

destroy:
	terraform destroy -auto-approve $(TF_OPTIONS)

mosh: _ssh_config
	mosh --ssh "ssh -F _ssh_config" test-machine

serve: _ssh_config _provision.ok
	ssh -F _ssh_config test-machine ./rtun-server

_init.ok: *.tf
	terraform init
	@touch $@

_server.ok: _init.ok
	terraform apply -auto-approve $(TF_OPTIONS)
	@touch $@

_provision.ok: provision.yml _hosts rtun rtun.yml rtun-server rtun-server.yml
	ansible-playbook provision.yml
	@touch $@

_ssh_config: files/ssh_config.in _server.ok
	terraform output ssh_config > $@

_hosts: files/hosts.in _server.ok
	terraform output hosts > $@

rtun.yml: files/rtun.yml.in _server.ok
	terraform output rtun > $@

rtun-server.yml: files/rtun-server.yml.in _server.ok
	terraform output rtun_server > $@

rtun:
	go build -o $@ ../agent/cmd

rtun-server:
	GOOS=linux GOARCH=amd64 go build -o $@ ../server/cmd
