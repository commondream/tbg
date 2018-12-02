# tbg - Trunk Based Git

A few commands that make it easier to work from trunk locally, even when
your team doesn't.

## Commands

### git unmerged

`git unmerged` lists commits on `master` that aren't on `origin/master`

It doesn't take any arguments and is a really simple wrapper for running `git logmaster ^origin/master --no-merges --pretty=oneline --abbrev-commit`.

### git share [branch] [rev]

`git share` pushes a particular commit to a remote branch for review

It does take arguments - the branch to push to and the commits to add. For example, running:

`git share new-feature HEAD` would add the `HEAD` commit to a new remote branch named `new-feature`.

`git share` does take all of the interesting shortcuts in commit names that git has to offer (like `HEAD`) by using `git rev-parse` to figure out the commit argument.

The command also handles subsequent runs by appending subsequent revs to the branch name given. This is handy when you get feedback during a review and need to modify your changes.

## Installation

```
go get -u github.com/commondream/tbg/...
```

## License

Copyright (c) 2018, Alan Johnson
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of the copyright holder nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
