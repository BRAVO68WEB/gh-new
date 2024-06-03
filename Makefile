install:
	@echo "> Installing wclone..."
	gh extension install .

remove:
	@echo "< Removing wclone..."
	gh extension remove wclone

reload: remove install