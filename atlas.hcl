data "external_schema" "gorm" {
  program = [
    "go", "run", "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./internal/models",
    "--dialect", "mysql",
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url

  # Atlas spins up a clean temporary MySQL via Docker to compute the diff.
  # If you don't have Docker, replace this with a local empty DB:
  # dev = "mysql://root:yourpassword@localhost:3306/countries_api_dev"
  dev = "docker://mariadb/latest/dev"

  migration {
    dir = "file://migrations"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
