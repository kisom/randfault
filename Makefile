TARGET :=	randfault
USER :=		kisom
REPO :=
VERSION :=	1.0.3
SOURCES :=	$(shell find . -name \*.go)

.PHONY: all
all: $(TARGET)

$(TARGET): $(SOURCES)
	CGO_ENABLED=0 go build

.PHONY: container
container: $(TARGET)
	docker build --rm -t $(USER)/$(TARGET):$(VERSION) .

.PHONY: publish
publish:
	docker push $(REPO)$(USER)/$(TARGET):$(VERSION)

.PHONY: clean
clean:
	rm -f $(TARGET)

print-%: ; @echo $*=$($*)
