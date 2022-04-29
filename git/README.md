vim技巧


## 如何避免不必要的commit
### 不修改commit message
- git commit --amend --no-edit
### 修改commit message
- git commit --amend


## 时光机
- git reflog
- git reset HEAD@{3} --hard


## 后悔药
### 1
- git log (找到前一个版本)
- git checkout ff33df18ebcaec423c152e3422d90aa4137fb310 -- README.md
- git commit / git reset --hard(取消后悔药操作)

### 2
- git revert f725945941fb0f996f4dd756ba198b33cd1fe06b
