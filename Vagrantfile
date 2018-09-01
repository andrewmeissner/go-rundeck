# standup with the --no-parallel flag

numRundecks = 2

Vagrant.configure("2") do |config|
    #  /var/rundeck/projects/<PROJECT>/etc/resources.xml

    ENV["VAGRANT_NO_PARALLEL"] = "true"

    config.vm.define "postgres-rundeck" do |postgres|
        postgres.vm.provider "docker" do |d|
            d.image = "postgres:10.5-alpine"
            d.name = "postgres-rundeck"
            d.ports = ["5432:5432"]
            d.env = {
                "POSTGRES_PASSWORD": "rundeckpassword",
                "POSTGRES_USER": "rundeck",
                "POSTGRES_DB": "rundeckdb"
            }
        end
    end

    (1..numRundecks).each do |i|
        config.vm.define "rundeck-#{i}" do |rundeck|
            rundeck.vm.provider "docker" do |d|
                d.image = "rundeck/rundeck:3.0.5"
                d.name = "rundeck-#{i}"
                d.ports = ["444#{i-1}:444#{i-1}"]
                d.cmd = ["-Dserver.http.port=444#{i-1}"]
                d.env = {
                    "RUNDECK_DATABASE_DRIVER": "org.postgresql.Driver",
                    "RUNDECK_DATABASE_USERNAME": "rundeck",
                    "RUNDECK_DATABASE_PASSWORD": "rundeckpassword",
                    "RUNDECK_DATABASE_URL": "jdbc:postgresql://postgres:5432/rundeckdb",
                    "RUNDECK_GRAILS_URL": "http://localhost:444#{i-1}",
                    "RUNDECK_STORAGE_PROVIDER": "db"
                }
                d.link("postgres-rundeck:postgres")
            end
        end
    end
end