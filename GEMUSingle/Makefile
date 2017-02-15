BINARY=GEMUSingle

.DEFAULT_GOAL: $(BINARY)
.PHONY: ${BINARY}

$(BINARY):
	go build -o ${BINARY}

.PHONY: install
install:
	go install ./...

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
