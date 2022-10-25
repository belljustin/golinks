build:
	sam build -t template.yaml -s ../../.

deploy: build
	sam deploy --config-file samconfig.toml
