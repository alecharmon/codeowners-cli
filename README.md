### Codeowners CLI tool 
Find files that are not covered by a repo's code owner file. Built on [codeowners](http://github.com/alecharmon/codeowners), a complete CODEOWNERS parser and blazing fast lookup tool written in Go.



- [Codeowners CLI tool](#codeowners-cli-tool)
  - [Install](#install)
  - [Usage](#usage)
    - [Check Owners](#check-owners)
    - [Verify CODEOWNERS File](#verify-codeowners-file)
    - [Help](#help)
#### Install
You can install this client globally via :
`go get -u github.com/alecharmon/codeowners-cli`

#### Usage
##### Check Owners
The root command lets you You can run the tool to check between any two commits and see if the related file changes had codeowners.
`codeowners-cli --base={base_commit} --head={commit_to_use_as_head}`
here is an example: 
![check](/images/check.png)
##### Verify CODEOWNERS File
Its also good to know that your CODEOWNERS file is valid as well, you can do that via the `verify` command.
![verify](/images/verify.png)

##### Help
Access more info via help flag, `--help`.
