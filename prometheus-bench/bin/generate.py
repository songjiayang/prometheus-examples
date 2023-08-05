from jinja2 import Template

def generate_prometheus_config(count):
    f = open("./config/prometheus.yml.template", "r")
    template = Template(f.read())
    f.close()

    wf = open("./config/prometheus.yml", "w")
    wf.write(template.render(count=count))
    wf.close()

def generate_hosts(count):
    f = open("./config/hosts.template", "r")
    template = Template(f.read())
    f.close()

    wf = open("./config/hosts", "w")
    wf.write(template.render(count=count))
    wf.close()


def main():
   total_count = 3400
   generate_prometheus_config(total_count)
   generate_hosts(total_count)
    
    
if __name__== "__main__" :
    main()