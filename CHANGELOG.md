# Changelog

## [1.8.2](https://github.com/soerenschneider/gobot-lux/compare/v1.8.1...v1.8.2) (2024-01-26)


### Bug Fixes

* **deps:** bump github.com/go-playground/validator/v10 ([#50](https://github.com/soerenschneider/gobot-lux/issues/50)) ([c38faea](https://github.com/soerenschneider/gobot-lux/commit/c38faea79390a92c753bc8dc3944ab355baf2337))
* **deps:** bump github.com/prometheus/client_golang ([#48](https://github.com/soerenschneider/gobot-lux/issues/48)) ([ed16c9f](https://github.com/soerenschneider/gobot-lux/commit/ed16c9f1c2dcfc508553b0864931d50ef338ac7a))
* **deps:** bump gobot.io/x/gobot/v2 from 2.1.1 to 2.3.0 ([#49](https://github.com/soerenschneider/gobot-lux/issues/49)) ([0d38b11](https://github.com/soerenschneider/gobot-lux/commit/0d38b11045c3939fdd855fd295b4c8e65e40b480))
* fix potential data race ([1eb64eb](https://github.com/soerenschneider/gobot-lux/commit/1eb64ebc121be9af0d8908c2e956c995b3afe6a8))

## [1.8.1](https://github.com/soerenschneider/gobot-lux/compare/v1.8.0...v1.8.1) (2024-01-26)


### Bug Fixes

* **deps:** bump github.com/go-playground/validator/v10 ([#39](https://github.com/soerenschneider/gobot-lux/issues/39)) ([498fa2f](https://github.com/soerenschneider/gobot-lux/commit/498fa2f822fe8810d6e8baaff970955e20256d73))
* **deps:** bump github.com/prometheus/client_golang ([#35](https://github.com/soerenschneider/gobot-lux/issues/35)) ([3867c41](https://github.com/soerenschneider/gobot-lux/commit/3867c412d3e93c68a17f1990eaff49738925195c))
* fix deadlock ([103c2ba](https://github.com/soerenschneider/gobot-lux/commit/103c2ba95a70d1bcd05fbb3ee835068696042738))

## [1.8.0](https://github.com/soerenschneider/gobot-lux/compare/v1.7.1...v1.8.0) (2023-07-14)


### Features

* send sensor reading instantly if exceeds deviation threshold ([752194a](https://github.com/soerenschneider/gobot-lux/commit/752194ad8259a9eb92f2c486002e8e47a21369d0))


### Bug Fixes

* fix potential race condition ([9b48682](https://github.com/soerenschneider/gobot-lux/commit/9b48682ab2798cc943d00b04d92fd60110c97b2a))
* fix struct field tags ([ba49dd6](https://github.com/soerenschneider/gobot-lux/commit/ba49dd687c33be68f7fc53ed2cc895b207f3e011))
* fix validation for tcp client certs ([a52243f](https://github.com/soerenschneider/gobot-lux/commit/a52243f5c1709c8aaf2b1e5e12ab2fcb0187650a))
* make server ca independent of crt and key ([26a6003](https://github.com/soerenschneider/gobot-lux/commit/26a6003f525c60eb4d94d10a5b87b31220e8c151))

## [1.7.1](https://github.com/soerenschneider/gobot-lux/compare/v1.7.0...v1.7.1) (2022-11-29)


### Miscellaneous Chores

* release 1.7.1 ([c03aad2](https://github.com/soerenschneider/gobot-lux/commit/c03aad2d74be2df531daa8e5198d409d4808be5f))

## [1.7.0](https://www.github.com/soerenschneider/gobot-lux/compare/v1.6.0...v1.7.0) (2022-05-04)


### Features

* enable tls client cert auth ([f5e43b9](https://www.github.com/soerenschneider/gobot-lux/commit/f5e43b9b5b82ad1cc884668daee56241e9c3cd5e))

## [1.6.0](https://www.github.com/soerenschneider/gobot-lux/compare/v1.5.1...v1.6.0) (2021-12-08)


### Features

* calculate avg in interval statistics ([0434246](https://www.github.com/soerenschneider/gobot-lux/commit/0434246ac7d0ff9d2c08a99d88ce1b99dc13237d))

### [1.5.1](https://www.github.com/soerenschneider/gobot-lux/compare/v1.5.0...v1.5.1) (2021-11-22)


### Bug Fixes

* add missing label ([7d1c287](https://www.github.com/soerenschneider/gobot-lux/commit/7d1c28747a4433c51c848cb43b33616dd0a4e11c))
* set to current timestamp instead of increasing ([92c07d0](https://www.github.com/soerenschneider/gobot-lux/commit/92c07d0e36891f307f0e463008ed5bd5630a3358))

## [1.5.0](https://www.github.com/soerenschneider/gobot-lux/compare/v1.4.0...v1.5.0) (2021-11-22)


### Features

* Collect statistics over configurable intervals ([8c4593c](https://www.github.com/soerenschneider/gobot-lux/commit/8c4593c2bebcda2e97ffc50e79a78fbcb54672be))

## [1.4.0](https://www.github.com/soerenschneider/gobot-lux/compare/v1.3.1...v1.4.0) (2021-11-02)


### Features

* add metric heartbeat ([c7a3d11](https://www.github.com/soerenschneider/gobot-lux/commit/c7a3d11588fa561b0b0931db8cd990afbb450d19))

### [1.3.1](https://www.github.com/soerenschneider/gobot-lux/compare/v1.3.0...v1.3.1) (2021-11-02)


### Miscellaneous Chores

* Trigger release ([114216f](https://www.github.com/soerenschneider/gobot-lux/commit/114216fbcd9dfe916d9a8b0580b2d3bedce93e46))

## [1.3.0](https://www.github.com/soerenschneider/gobot-lux/compare/v1.2.0...v1.3.0) (2021-10-21)


### Features

* Add version info metric ([9c0c23b](https://www.github.com/soerenschneider/gobot-lux/commit/9c0c23b0dd120e7cbb5eac81c40b9d7a42dd3514))


### Bug Fixes

* more reasonable limit for polling interval ([54652a0](https://www.github.com/soerenschneider/gobot-lux/commit/54652a01a066f0daf86ff56fb86ea2b3d5b3a289))
* set auto-reconnect to true ([a8963a3](https://www.github.com/soerenschneider/gobot-lux/commit/a8963a39570f6f46b072b17bd6b64e2ea890f8dd))

## [1.2.0](https://www.github.com/soerenschneider/gobot-lux/compare/v1.1.0...v1.2.0) (2021-10-21)


### Features

* add flag to print version info ([820f66d](https://www.github.com/soerenschneider/gobot-lux/commit/820f66d90871217a51f6c77ec4f32ca57e96bc44))
* reproducible builds by omitting build time ([c6c60f9](https://www.github.com/soerenschneider/gobot-lux/commit/c6c60f9d22c7f677f8f10f4dcaaeb4e444c7d1ef))


### Bug Fixes

* use provided wrapper to send msg ([c7b241f](https://www.github.com/soerenschneider/gobot-lux/commit/c7b241f20cda5ff7f9aa4634d146265f577a0243))

## [1.1.0](https://www.github.com/soerenschneider/gobot-brightness/compare/v1.0.1...v1.1.0) (2021-09-14)


### Features

* print version ([e5ff4a1](https://www.github.com/soerenschneider/gobot-brightness/commit/e5ff4a1044c974363964cf254f7069b0488d08a9))

### [1.0.1](https://www.github.com/soerenschneider/gobot-brightness/compare/v1.0.0...v1.0.1) (2021-09-14)


### Bug Fixes

* Visibility of function ([521e26d](https://www.github.com/soerenschneider/gobot-brightness/commit/521e26d7dfc06726ae6c35a17ed22da5b3f64784))

## 1.0.0 (2021-09-13)


### Miscellaneous Chores

* release 1.0.0 ([790c02c](https://www.github.com/soerenschneider/gobot-brightness/commit/790c02c012c5eef52f64ab68dfafb14c9cb828b6))
