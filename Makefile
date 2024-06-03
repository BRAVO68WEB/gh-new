install:
	@echo "> Installing new..."
	gh extension install .

remove:
	@echo "< Removing new..."
	gh extension remove new

reload: remove install