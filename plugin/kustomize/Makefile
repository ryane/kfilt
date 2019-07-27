REGISTRY := ryane

.PHONY: build
build:
	@docker build -t ${REGISTRY}/kfilt-kustomize-plugin .

.PHONY: run
run:
	@docker run --rm -it -v $(PWD):/bases \
		${REGISTRY}/kfilt-kustomize-plugin:latest \
		build --enable_alpha_plugins /bases/example
