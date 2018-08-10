# standup with the --no-parallel flag

Vagrant.configure("2") do |config|
    #  /var/rundeck/projects/<PROJECT>/etc/resources.xml

    ENV["VAGRANT_NO_PARALLEL"] = "true"

    config.vm.define "postgres" do |postgres|
        postgres.vm.provider "docker" do |d|
            d.image = "postgres:10.4-alpine"
            d.name = "postgres"
            d.ports = ["5432:5432"]
            d.env = {
                "POSTGRES_PASSWORD": "rundeckpassword",
                "POSTGRES_USER": "rundeck",
                "POSTGRES_DB": "rundeckdb"
            }
        end
    end

    rundeck_versions = ["2.11.5", "2.11.5"]

    (1..rundeck_versions.length).each do |i|
        config.vm.define "rundeck-#{i}" do |rundeck|
            rundeck.vm.provider "docker" do |d|
                d.image = "jordan/rundeck:#{rundeck_versions[i-1]}"
                d.name = "rundeck-#{i}"
                d.ports = ["444#{i-1}:4440"]
                d.env = {
                    "RUNDECK_PASSWORD": "rundeckpassword",
                    "RUNDECK_ADMIN_PASSWORD": "admin",
                    "EXTERNAL_SERVER_URL": "http://localhost:444#{i-1}",
                    "CLUSTER_MODE": "true",
                    "DATABASE_URL": "jdbc:postgresql://postgres:5432/rundeckdb",
                    "RUNDECK_STORAGE_PROVIDER": "db",
                    "RUNDECK_PROJECT_STORAGE_TYPE": "db",
                    "NO_LOCAL_MYSQL": "true"
                }
                d.link("postgres:postgres")
            end
        end
    end    
end