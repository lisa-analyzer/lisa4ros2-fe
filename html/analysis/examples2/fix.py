# STEP 1:
# For every project
#  - fix report.html and extract info

import os
import re
import pandas

repos = {}
#file = pandas.read_csv(env.PARSER_FILE_PATH)
#file = pandas.read_csv("lisaros2-analysis-02032024-results-single.csv", sep=";")
#file = pandas.read_csv("lisaros2-analysis-04032024-results-single-filter-EMPTY-TOPBOTTOM.csv", sep=";")
file = pandas.read_csv("repos.csv", sep=";")
for idx, line in file.iterrows():
    #path = line["file_name"]
    proj = line["project"]
    url = line["github_url"]
    repos[proj] = url
projects = {}
entry = """<li>
          <h3>{name}</h3>
          <p>(Repository: <a href="{github}">{github}</a>)</p>
            <ul>
            <li>
                <i>(<a href="{name}/report.html">View Analysis</a> | <a href="{name}/sources.zip">Download sources</a>)</i>
                {results}
            </li>
            </ul>
            </li>
            """
for item in os.listdir("."):
    if os.path.isdir(item):
        with open(item + "/report.html", "r") as f:
            lines = "".join(f.readlines())
            lines = lines.replace("<div class=\"ml-4\"><h4><a href=\"/\">Home</a></h4></div>", "<div class=\"ml-4\"><h4><a href=\"../index.html\">Home</a></h4></div>")
            #lines = lines.replace("href=\"..\"", "href=\"../index.html\"")
            regex = re.findall(r"<div>\s*<h2>Analysis Results<\/h2>[\s\S]*<\/ul>\s*<\/div>", lines)
            projects[item] = regex[0].replace("<h2>Analysis Results</h2>", "")
        with open(item + "/report.html", "w") as f:
            f.write(lines)
print(repos)
entries = ""
for project in projects:
    entries += entry.format(name=project, github=repos[project], results=projects[project]) + "\n"
#print(projects)
template = ""
with open("index_template.html", "r") as f:
    template = "".join(f.readlines())
    template = template.replace("{examples}", entries)
with open("index.html", "w") as f:
    f.write(template)