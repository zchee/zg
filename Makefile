#
#  Makefile for Go
#
GO_CMD=go

GO_GCFLAGS :=
# Parse current git-commit hash, and overwride `GitCommit` string variable
# -ldflags '-X': Set the value of the string variable in importpath named name to value
GO_LDFLAGS := -X `go list `.GitCommit=`git rev-parse --short HEAD 2>/dev/null`

# Set gcflag, ldflags for gc compiler
# Usage: DEBUG=true make
ifeq ($(DEBUG),true)
	# Disable function inlining and variable registerization. For lldb, gdb, dlv and the involved debugger tools
	# See also Dave cheney's blog post: http://goo.gl/6QCJMj
	# And, My cgo blog post: http://libraryofalexandria.io/cgo/
	#
	# -gcflags '-N': Will be disable the optimisation pass in the compiler
	# -gcflags '-l': Will be disable inlining (but still retain other compiler optimisations)
	#                This is very useful if you are investigating small methods, but can’t find them in `objdump`
	GO_GCFLAGS := -gcflags "-N -l"
else
	# Turn of DWARF debugging information and strip the binary otherwise
	# It will reduce the as much as possible size of the binary
	# See also Russ Cox's answered in StackOverflow: http://goo.gl/vOaigc
	#
	# -ldflags '-w': Turns off DWARF debugging infomation
	# 	- Will not be able to use lldb, gdb, objdump or related to debugger tools
	# -ldflags '-s': Turns off generation of the Go symbol table
	# 	- Will not be able to use `go tool nm` to list symbols in the binary
	# 	- `strip -s` is like passing '-s' flag to -ldflags, but it doesn't strip quite as much
	GO_LDFLAGS := $(GO_LDFLAGS) -w -s
endif

# Set static build option
# Usage: STATIC=true make
ifeq ($(STATIC),true)
	# Append to the version
	GO_LDFLAGS := $(GO_LDFLAGS) -extldflags -static
endif

# Global go command environment variables
GO_BUILD=$(GO_CMD) build -o $(OUTPUT_NAME)
GO_BUILD_RACE=$(GO_CMD) build -race -o $(OUTPUT_NAME)
GO_TEST=$(GO_CMD) test
GO_TEST_VERBOSE=$(GO_CMD) test -v
GO_RUN=$(GO_CMD) run
GO_INSTALL=$(GO_CMD) install -v
GO_CLEAN=$(GO_CMD) clean
GO_DEPS=$(GO_CMD) get -d -v
GO_DEPS_UPDATE=$(GO_CMD) get -d -v -u
GO_VET=$(GO_CMD) vet
GO_FMT=$(GO_CMD) fmt
GO_LINT=golint

# Color output
CRESET=\x1b[0m
CRED=\x1b[31;01m
CGREEN=\x1b[32;01m
CYELLOW=\x1b[33;01m
CBLUE=\x1b[34;01m
CMAGENTA=\x1b[35;01m
CCYAN=\x1b[36;01m

# Package infomation
GITHUB_USER=zchee
TOP_PACKAGE_DIR := github.com/$(GITHUB_USER)
PACKAGE_LIST := `basename $(PWD)`
OUTPUT_NAME := `basename $(PWD)`
# Parse "func main()" only '.go' file on current dir
# FIXME: Not support main.go
MAINFILE_NAME := `grep "func main\(\)" *.go -l`

# gotags
# go get -u -v github.com/jstemmer/gotags
CTAGS_CMD=gotags
# gotags options
#  	-f string
# 		Output specified file name
# 		If specified "-", output to stdout
# 	-R
# 		Recurse into directories in the file list
# 	-fields string
# 		Include selected extension fields (only +l)
# 	-sort
# 		Sort tags (default true)
# 	-tag-relative
# 		File path s should be relative to the directory containing the tag file
CTAGS_OPTIONS=-f tags -R -fields=+l -sort -tag-relative .

all: build

build: vet
	@for p in $(PACKAGE_LIST); do \
		echo "$(CBLUE)==>$(CRESET) Build $(CGREEN)$$p$(CRESET) ..."; \
		$(GO_BUILD) -ldflags "$(GO_LDFLAGS)" $(GO_GCFLAGS) $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

build-race: vet
	@for p in $(PACKAGE_LIST); do \
		echo "$(CBLUE)==>$(CRESET) Build $(CGREEN)$$p$(CRESET) with -race flag ..."; \
		$(GO_BUILD_RACE) -ldflags "$(GO_LDFLAGS)" $(GO_GCFLAGS) $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

build-force: vet
	@for p in $(PACKAGE_LIST); do \
		echo "$(CBLUE)==>$(CRESET) Build $(CGREEN)$$p$(CRESET) with force rebuilding of packages ..."; \
		$(GO_BUILD) -a -ldflags "$(GO_LDFLAGS)" $(GO_GCFLAGS) $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

build-verbose: vet
	@for p in $(PACKAGE_LIST); do \
		echo "$(CBLUE)==>$(CRESET) Build $(CGREEN)$$p$(CRESET) with verbose ..."; \
		$(GO_BUILD) -v -x -ldflags "$(GO_LDFLAGS)" $(GO_GCFLAGS) $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

run: vet
	@for p in $(PACKAGE_LIST); do \
		echo "$(CBLUE)==>$(CRESET) Run $(CGREEN)$$p$(CRESET) ..."; \
		$(GO_RUN) $(MAINFILE_NAME) || exit 1; \
	done

test: deps
	@for p in $(PACKAGE_LIST); do \
		echo "$(CBLUE)==>$(CRESET) Unit Testing $(CGREEN)$$p$(CRESET) ..."; \
		$(GO_TEST) $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

test-verbose: deps
	@for p in $(PACKAGE_LIST); do \
		echo "$(CBLUE)==>$(CRESET) Unit Testing $(CGREEN)$$p$(CRESET) ..."; \
		$(GO_TEST_VERBOSE) $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

deps:
	@for p in $(PACKAGE_LIST); do \
		echo "$(CBLUE)==>$(CRESET) Install dependencies for $(CGREEN)$$p$(CRESET) ..."; \
		$(GO_DEPS) $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

update-deps:
	@for p in $(PACKAGE_LIST); do \
		echo "$(CBLUE)==>$(CRESET) Update dependencies for $(CGREEN)$$p$(CRESET) ..."; \
		$(GO_DEPS_UPDATE) $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

install:
	@for p in $(PACKAGE_LIST); do \
		echo "$(CBLUE)==>$(CRESET) Install $(CGREEN)$$p$(CRESET) ..."; \
		$(GO_INSTALL) $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

clean:
	@for p in $(PACKAGE_LIST); do \
		echo "$(CBLUE)==>$(CRESET) Clean $(CGREEN)$$p$(CRESET) ..."; \
		$(GO_CLEAN) $(TOP_PACKAGE_DIR)/$$p; \
	done

fmt:
	@for p in $(PACKAGE_LIST); do \
		echo "$(CBLUE)==>$(CRESET) Formatting $(CGREEN)$$p$(CRESET) ..."; \
		$(GO_FMT) $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

vet:
	@for p in $(PACKAGE_LIST); do \
		echo "$(CBLUE)==>$(CRESET) Vet $(CGREEN)$$p$(CRESET) ..."; \
		$(GO_VET) $(TOP_PACKAGE_DIR)/$$p; \
	done

lint:
	@for p in $(PACKAGE_LIST); do \
		echo "$(CBLUE)==>$(CRESET) Lint $(CGREEN)$$p$(CRESET) ..."; \
		$(GO_LINT) src/$(TOP_PACKAGE_DIR)/$$p; \
	done

ctags:
	@for p in $(PACKAGE_LIST); do \
		echo "$(CBLUE)==>$(CRESET) Create ctags file..."; \
		$(CTAGS_CMD) $(CTAGS_OPTIONS) || exit 1; \
	done

make-debug:
	# defined command list
	@echo GO_CMD=$(GO_CMD)
	@echo GO_BUILD=$(GO_BUILD)
	@echo GO_BUILD_RACE=$(GO_BUILD_RACE)
	@echo GO_TEST=$(GO_TEST)
	@echo GO_TEST_VERBOSE=$(GO_TEST_VERBOSE)
	@echo GO_INSTALL=$(GO_INSTALL)
	@echo GO_CLEAN=$(GO_CLEAN)
	@echo GO_DEPS=$(GO_DEPS)
	@echo GO_DEPS_UPDATE=$(GO_DEPS_UPDATE)
	@echo GO_VET=$(GO_VET)
	@echo GO_FMT=$(GO_FMT)
	@echo GO_LINT=$(GO_LINT)
	@echo TOP_PACKAGE_DIR=$(TOP_PACKAGE_DIR)
	@echo PACKAGE_LIST=$(PACKAGE_LIST)
	@echo MAINFILE_NAME=$(MAINFILE_NAME)
	@echo OUTPUT_NAME=$(OUTPUT_NAME)

.PHONY: all build build-race test test-verbose deps update-deps install clean fmt vet lint
