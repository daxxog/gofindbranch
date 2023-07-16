# 1
make && ./gofindbranch -filter '.+' -current-user-only | tee repos.txt
# 2
make && cat repos.txt | cat repos.txt | sed 's/^/daxxog\//g' | ./gofindbranch -filter '^dv.+'
# 3
make && cat repos.txt | cat repos.txt | head | sed 's/^/daxxog\//g' | ./gofindbranch -filter '.+' -open-prs-only
# 4
make && ./gofindbranch repos -f '.+' --current-user-only | head | sed 's/^/daxxog\//g' | ./gofindbranch branches -f '.+' --open-prs-only
