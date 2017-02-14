.DEFAULT_GOAL: GEMUSingle

.PHONY: GEMUSingle
GEMUSingle:
	cd GEMUSingle && make
.PHONY: GEMUSingle_clean
GEMUSingle_clean:
	cd GEMUSingle && make clean
.PHONY: GEMUSingle_install
GEMUSingle_install:
	cd GEMUSingle && make install

.PHONY: clean
clean: GEMUSingle_clean
.PHONY: install
install: GEMUSingle_install
