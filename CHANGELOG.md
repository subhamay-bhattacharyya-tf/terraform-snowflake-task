# Changelog

All notable changes to this project will be documented in this file.

## [2.0.0](https://github.com/subhamay-bhattacharyya-tf/terraform-snowflake-task/compare/v1.0.0...v2.0.0) (2026-03-02)

### ⚠ BREAKING CHANGES

* Remove modules/ directory structure

- Move module files to repository root
- Update examples to reference root module
- Update CI workflow paths
- Align with one-module-per-repo best practice

### Features

* convert to single-module repository layout ([dc33b84](https://github.com/subhamay-bhattacharyya-tf/terraform-snowflake-task/commit/dc33b847723773f97561aa625f3e6fd2475b28c2))

### Bug Fixes

* **snowflake-task:** refine schedule block conditional logic ([9f7393c](https://github.com/subhamay-bhattacharyya-tf/terraform-snowflake-task/commit/9f7393c81da913012cfe29c7d23571e4c2995122))

## 1.0.0 (2026-02-09)

### Features

* implement Snowflake task module ([2b6f870](https://github.com/subhamay-bhattacharyya-tf/terraform-snowflake-task/commit/2b6f870ba7ed5938f22a53f10dd15e35daff9966))

### Bug Fixes

* **snowflake-task:** correct warehouse size validation logic ([69604da](https://github.com/subhamay-bhattacharyya-tf/terraform-snowflake-task/commit/69604da426f79d5b2ab0a077766a45d867dc459f))
* **snowflake-task:** handle empty task dependencies gracefully ([285e094](https://github.com/subhamay-bhattacharyya-tf/terraform-snowflake-task/commit/285e094ad00e5955731e6eb5d54752eec117e089))

## [unreleased]

### 🚀 Features

- Implement Snowflake task module
- [**breaking**] Convert to single-module repository layout

### 🐛 Bug Fixes

- *(snowflake-task)* Correct warehouse size validation logic
- *(snowflake-task)* Handle empty task dependencies gracefully
- *(snowflake-task)* Refine schedule block conditional logic

### 📚 Documentation

- Update CHANGELOG.md [skip ci]
- Update CHANGELOG.md [skip ci]
- Update CHANGELOG.md [skip ci]
- Update CHANGELOG.md [skip ci]
- Update CHANGELOG.md [skip ci]
- Update CHANGELOG.md [skip ci]
- Update CHANGELOG.md [skip ci]
- Update task dependency references to use fully qualified names
- Update CHANGELOG.md [skip ci]
- *(examples)* Update task dependency references with quoted identifiers
- Update CHANGELOG.md [skip ci]
- *(examples)* Remove quoted identifiers from task dependency references
- Update CHANGELOG.md [skip ci]
- Update CHANGELOG.md [skip ci]
- Update CHANGELOG.md [skip ci]

### 🧪 Testing

- *(snowflake-task)* Update test assertions for type consistency
- *(multiple_tasks)* Remove transform and load task test cases
- *(multiple_tasks)* Remove unused task name variables

### ⚙️ Miscellaneous Tasks

- Upgrade Snowflake provider to 1.0.0 and improve documentation
