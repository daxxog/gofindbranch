# 1
make && ./gofindbranch -filter '.+' -current-user-only | tee repos.txt
# 2
make && cat repos.txt | cat repos.txt | sed 's/^/daxxog\//g' | ./gofindbranch -filter '^dv.+'
