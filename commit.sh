msg=$1
type=$2

if [[ -z "${msg}" ]]
then
	mgs="fix: Bug fixes and updates."
fi

if [[ -z "${type}" ]]
then
	type="fix"
fi

# commit
git add -A .
git commit -m "${type}: ${msg}"
git push -u origin dev
