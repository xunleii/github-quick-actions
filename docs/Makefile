# Copyright (c) 2021 Alexandre NICOLAIE
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.
# -----------------------------------------------------------------------------

# This Makefile provides common helpers to manage documentations and
# repository administration, like creating issues on new Quick Actions.

DOCS_DIR	:=	$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

# -----------------------------------------------------------------------------
# Manage README.md generation

README_TPL_PATH			?=	${DOCS_DIR}/README.tpl.md
QUICK_ACTIONS_LIST_PATH	?=	${DOCS_DIR}/quick_actions.list.toml

README.md: README.md@pt.01 README.md@pt.02 README.md@pt.03

README.md@pt.01:
	@sed '/::include quick_actions_table/Q' ${README_TPL_PATH}
README.md@pt.02: README.md/tables/released_quick_actions README.md/tables/todo_quick_actions README.md/tables/rejected_quick_actions
README.md@pt.03:
	@sed -e '1,/::include quick_actions_table/ d' ${README_TPL_PATH}

README.md/tables/released_quick_actions:
	@echo "## Available quick actions"
	@echo
	@echo "The following quick actions are already released and available on the Github application."
	@echo
	@echo "|     Command     |   Applicable on   |  Description  |"
	@echo "| :-------------: | :---------------- | :-----------: |"
	@tomlq '.' ${QUICK_ACTIONS_LIST_PATH} | jq '.quick_actions.released[] | "| \(.quick_action | map("`\(.)`") | join("<br>")) | \(.on_events | map("**&#10003;** `\(.)`") | join("<br>")) | \(.description |= gsub("\n"; "<br>") | .description) |"' -r
	@echo

README.md/tables/todo_quick_actions:
	@echo "## Quick actions to be developed"
	@echo
	@echo "The following quick actions will be available in the future (must need times to develop them)."
	@echo
	@echo "|     Command     |   Applicable on   |  Description  |"
	@echo "| :-------------: | :---------------- | :-----------: |"
	@tomlq '.' ${QUICK_ACTIONS_LIST_PATH} | jq '.quick_actions.next_releases[] | "| \(.quick_action | map("`\(.)`") | join("<br>")) | \(.on_events | map("**&#9676;** `\(.)`") | join("<br>")) | \(.description |= gsub("\n"; "<br>") | .description) |"' -r
	@echo

README.md/tables/rejected_quick_actions:
	@echo "## Quick actions that will not be developed"
	@echo
	@echo "The following quick actions will not be developed for specific reasons."
	@echo
	@echo "|     Command     | Description |  Reasons  |"
	@echo "| :-------------: | :---------- | :-------: |"
	@tomlq '.' ${QUICK_ACTIONS_LIST_PATH} | jq '.quick_actions.rejected[] | "| \(.quick_action | map("`\(.)`") | join("<br>")) | \(.description |= gsub("\n"; "<br>") | .description) | \(.reasons | map(gsub("\n"; "<br>")) | join("<br>")) |"' -r
	@echo

