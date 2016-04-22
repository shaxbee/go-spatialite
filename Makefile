BUILDTAGS=-tags "libsqlite3"

test:
	go test -v . $(BUILDTAGS)
	go test -v ./wkb

cover:
	go test -v . -covermode=count -coverprofile=profile.cov $(BUILDTAGS)
	go test -v ./wkb -covermode=count -coverprofile=wkb/profile.cov 
	gocovmerge profile.cov wkb/profile.cov > merged.cov

coverhtml: cover
	go tool cover -html=merged.cov	

install:
	go install $(BUILDTAGS)
