| Authors                            | Creation Date | Status      | Extra |
|------------------------------------|---------------|-------------|---|
| @camilamacedo86,@Kavinjsir,@varshaprasad96 | Feb, 2023     | Implementable | - |

Experimental Helper to upgrade projects by re-scaffolding
===================

This proposal aims to provide a new alpha command with a helper which
would be able to re-scaffold the project from the scratch based on
the [PROJECT config][project-config].

## Example

By running a command like following, users would be able to re-scaffold the whole project from the scratch using the
current version of KubeBuilder binary available.

```shell
kubebuilder alpha generate [OPTIONS]
```

### Workflows

Following some examples of the workflows

**To update the project with minor changes provided**

See that for each KubeBuilder release the plugins versions used to scaffold
the projects might have bug fixes and new incremental features added to the
templates which will result in changes to the files that are generated by
the tool for new projects.

In this case, you used previously the tool to generate the project
and now would like to update your project with the latest changes
provided for the same plugin version. Therefore, you will need to:

- Download and install KubeBuilder binary ( latest / upper release )
- You will run the command in the root directory of your project: `kubebuilder alpha generate`
- Then, the command will remove the content of your local directory and re-scaffold the project from the scratch
- It will allow you to compare your local branch with the remote branch of your project to re-add the code on top OR
  if you do not use the flag `--no-backup` then you can compare the local directory with the copy of your project
  copied to the path `.backup/project-name/` before the re-scaffold be done.
- Therefore, you can run make all and test the final result. You will have after all your project updated.

**To update the project with major changes provided**

In this case, you are looking for to migrate the project from, for example,
`go/v3` to `go/v4`. The steps are very similar to the above ones. However,
in this case you need to inform the plugin that you want to use to do the scaffold
from scratch `kubebuilder alpha generate --plugins=go/v4`.

## Open Questions

N/A

## Summary

Therefore, a new command can be designed to load user configs from the [PROJECT config][project-config] file, and run the corresponding kubebuilder subcommands to generate the project based on the new kubebuilder version. Thus, it makes it easier for the users to migrate their operator projects to the new scaffolding.

## Motivation

A common scenario is to upgrade the project based on the newer Kubebuilder. The recommended (straightforward) steps are:

- a) re-scaffold all files from scratch using the upper version/plugins
- b) copy user-defined source code to the new layout

The proposed command will automate the process at maximum, therefore helping operator authors with minimizing the manual effort.

The main motivation of this proposal is to provide a helper for upgrades and
make less painful this process. Examples:

- See the discussion [How to regenerate scaffolding?](https://github.com/kubernetes-sigs/kubebuilder/discussions/2864)
- From [slack channel By Paul Laffitte](https://kubernetes.slack.com/archives/CAR30FCJZ/p1675166014762669)

### Goals

- Help users upgrade their project with the latest changes
- Help users to re-scaffold projects from scratch based on what was done previously with the tool
- Make less painful the process to upgrade

### Non-Goals

- Change the default layout or how the KubeBuilder CLI works
- Deal with customizations or deviations from the proposed layout
- Be able to perform the project upgrade to the latest changes without human interactions
- Deal and support external plugins
- Provide support to older version before having the Project config (Kubebuilder < 3x) and the go/v2 layout which exists to ensure  a backwards compatibility with legacy layout provided by Kubebuilder 2x

## Proposal

The proposed solution to achieve this goal is to create an alpha command as described
in the example section above, see:

```shell
kubebuilder alpha generate \
    --input-dir=<path where the PROJECT file can be found>
    --output-dir=<path where the project should be re-scaffold>
    --no-backup
    --backup-path=<path-where the current version of the project should be copied as backup>
    --plugins=<chain of plugins key that can be used to create the layout with init sub-command>
```

**Where**:

- input-dir: [Optional] If not informed, then, by default, it is the current directory (project directory). If the `PROJECT` file does not exist, it will fail.
- output-dir: [Optional] If not informed then, it should be the current repository.
- no-backup: [Optional] If not informed then, the current directory should be copied to the path `.backup/project-name`
- backup: [Optional] If not informed then, the backup will be copied to the path `.backup/project-name`
- plugins:  [Optional] If not informed then, it is the same plugin chain available in the layout field
- binary: [Optional] If not informed then, the command will use KubeBuilder binary installed globaly.

> Note that the backup created in the current directory must be prefixed with `.`. Otherwise the tool
will not able to perform the scaffold to create a new project from the scratch.

This command would mainly perform the following operations:

- 1. Check the flags
- 2. If the backup flag be used, then check if is a valid path and make a backup of the current project
- 3. Copy the whole current directory to `.backup/project-name`
- 4. Ensure that the output path is clean. By default it is the current directory project where the project was scaffolded previously and it should be cleaned up before to do the re-scaffold.
Only the content under `.backup/project-name` should be kept.
- 4. Read the [PROJECT config][project-config]
- 5. Re-run all commands using the KubeBuilder binary to recreate the project in the output directory

The command should also provide a comprensive help with examples of the proposed workflows. So that, users
are able to understand how to use it when run `--help`.

### User Stories

**As an Operator author:**

- I can re-generate my project from scratch based on the proposed helper, which executes all the
commands according to my previous input to the project. That way, I can easily migrate my project to the new layout
using the newer CLI/plugin versions, which support the latest changes, bug fixes, and features.
- I can regenerate my project from the scratch based on all commands that I used the tool to build
my project previously but informing a new init plugin chain, so that I could upgrade my current project to new
layout versions and experiment alpha ones.
- I would like to re-generate the project from the scratch using the same config provide in the PROJECT file and inform
a path to do a backup of my current directory so that I can also use the backup to compare with the new scaffold and add my custom code
on top again without the need to compare my local directory and new scaffold with any outside source.

**As a Kubebuiler maintainer:**

- I can leverage this helper to easily migrate tutorial projects of the Kubebuilder book.
- I can leverage on this helper to encourage its users to migrate to upper versions more often, making it easier to maintain the project.

### Implementation Details/Notes/Constraints

Note that in the [e2e tests](https://github.com/kubernetes-sigs/kubebuilder/tree/master/test/e2e) the binary is used to do the scaffolds.
Also, very similar to the implementation that exist in the integration test KubeBuilder has
a code implementation to re-generate the samples used in the docs and add customizations on top,
for further information check the [hack/docs](https://github.com/kubernetes-sigs/kubebuilder/tree/master/hack/docs).

This subcommand could have a similar implementation that could be used by the tests and this plugin.
Note that to run the commands using the binaries we are mainly using the following golang implementation:

```go
cmd := exec.Command(t.BinaryName, Options)
_, err := t.Run(cmd)
```

### Risks and Mitigations

**Hard to keep the command maintained**

A risk to consider is that it would be hard to keep this command maintained
because we need to develop specific code operations for each plugin. The mitigation for
this problem could be developing a design more generic that could work with all plugins.

However, initially a more generic design implementation does not appear to be achievable and
would be considered out of the scope of this proposal (no goal). It should to be considered
as a second phase of this implementation.

Therefore, the current achievable mitigation in place is that KubeBuilder's  policy of not providing official
support of maintaining and distributing many plugins.

### Proof of Concept

All input data is tracked. Also, as described above we have examples of code implementation
that uses the binary to scaffold the projects. Therefore, the goal of this project seems
very reasonable and achievable. An initial work to try to address this requirement can
be checked in this [pull request](https://github.com/kubernetes-sigs/kubebuilder/pull/3022)

## Drawbacks

- If the value that feature provides does not pay off the effort to keep it
  maintained, then we would need to deprecate and remove the feature in the long term.

## Alternatives

N/A

## Implementation History

The idea of automate the re-scaffold of the project is what motivates
us track all input data in to the [project config][project-config]
in the past. We also tracked the [issue](https://github.com/kubernetes-sigs/kubebuilder/issues/2068)
based on discussion that we have to indeed try to add further
specific implementations to do operations per major bumps. For example:

To upgrade from go/v3 to go/v4 we know exactly what are the changes in the layout
then, we could automate these specific operations as well. However, this first idea is harder yet
to be addressed and maintained.

## Future Vision

We could use it to do cool future features such as creating a GitHub action which would push-pull requests against the project repositories to help users be updated with, for example, minor changes. By using this command, we might able to git clone the project and to do a new scaffold and then use some [git strategy merge](https://www.geeksforgeeks.org/merge-strategies-in-git/) to result in a PR to purpose the required changes.

We probably need to store the CLI tool tag release used to do the scaffold to persuade this idea. So that we can know if the project requires updates or not.

[project-config]: https://book.kubebuilder.io/reference/project-config.html