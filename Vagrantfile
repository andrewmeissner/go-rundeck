Vagrant.configure("2") do |config|
   config.vm.define "rundeck" do |rundeck|
    rundeck.vm.provider "docker" do |d|
        d.image = "jordan/rundeck"
        d.name = "rundeck"
        d.ports = ["4440:4440"]
        d.env = {
            "RUNDECK_PASSWORD": "admin",
            "RUNDECK_ADMIN_PASSWORD": "admin",
            "EXTERNAL_SERVER_URL": "http://localhost:4440"
        }
    end
   end
end
