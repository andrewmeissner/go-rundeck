Vagrant.configure("2") do |config|
    #  /var/rundeck/projects/<PROJECT>/etc/resources.xml

    config.vm.define "postgres" do |postgres|
        postgres.vm.provider "docker" do |d|
            d.image = "postgres"
            d.name = "postgres"
            d.ports = ["5432:5432"]
            d.env = {
                "POSTGRES_PASSWORD": "rundeckpassword",
                "POSTGRES_USER": "rundeck",
                "POSTGRES_DB": "rundeckdb"
            }
        end
    end

    config.vm.define "rundeck-1" do |rundeck|
        rundeck.vm.provider "docker" do |d|
            d.image = "jordan/rundeck"
            d.name = "rundeck-1"
            d.ports = ["4440:4440"]
            d.env = {
                "RUNDECK_PASSWORD": "rundeckpassword",
                "RUNDECK_ADMIN_PASSWORD": "admin",
                "EXTERNAL_SERVER_URL": "http://localhost:4440",
                "CLUSTER_MODE": "true",
                "DATABASE_URL": "jdbc:postgresql://postgres:5432/rundeckdb",
                "RUNDECK_STORAGE_PROVIDER": "db",
                "RUNDECK_PROJECT_STORAGE_TYPE": "db",
                "NO_LOCAL_MYSQL": "true"
            }
            d.link("postgres:postgres")
        end
    end

    config.vm.define "rundeck-2" do |rundeck|
        rundeck.vm.provider "docker" do |d|
            d.image = "jordan/rundeck"
            d.name = "rundeck-2"
            d.ports = ["4441:4440"]
            d.env = {
                "RUNDECK_PASSWORD": "rundeckpassword",
                "RUNDECK_ADMIN_PASSWORD": "admin",
                "EXTERNAL_SERVER_URL": "http://localhost:4441",
                "CLUSTER_MODE": "true",
                "DATABASE_URL": "jdbc:postgresql://postgres:5432/rundeckdb",
                "RUNDECK_STORAGE_PROVIDER": "db",
                "RUNDECK_PROJECT_STORAGE_TYPE": "db",
                "NO_LOCAL_MYSQL": "true",
                "SKIP_DATABASE_SETUP": "true"
            }
            d.link("postgres:postgres")
        end
    end
end