# Configuration package

#### Пакет configuration(конфигурациии) содержит в себе все необходимые инструменты для взаимодействия с файлом конфигурации configuration.yaml.

---

Структура Configuration дает доступ ко всем полям, заданным в .env файле, посредством вызова функции ConfigurationInit.
> config, err := ConfigurationInit()


После этого появляется возможность обращаться к элементам структуры Configuration
> fmt.Println(config.Grafana.Address)
> > http://localhost:3000

### Поля структуры:

* Grafana - содержит адрес для обращения Address и API key.
* MetricsAPI - содержит url-адрес Address для обращения, в котором хранятся метрики.
* Schedule - это расписание для работы приложения, содержит поля Repeat, Days, Hours и Minutes.

---



#### Необходимость данного пакета обуславливается возможностью кастомизации готового решения.