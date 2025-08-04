print('Iniciando inicialização do MongoDB...');

// Conecta ao banco de dados appseclabs
db = db.getSiblingDB('appseclabs');

// Cria a collection labs
db.createCollection('labs');

// Insere os labs
db.labs.insertMany([
    {
        labSpec: {
            image: "ghcr.io/vitor-mauricio/lab-runner-example:dev",
            ports: [3000],
            code_config: {
                git_url: "github.com/appsec-digital/labs.git",
                git_branch: "main",
                git_path: "lab-example",
            },
            env: {
                "DOCKER_HOST": "tcp://localhost:2375",
                "APP_PORT": "3000",
                "APP_URL": "$(LAB_BASE_URL)/$(NAMESPACE)/proxy/$(APP_PORT)"
            },
            evaluations: [
                {
                    name: "sast",
                    order: 1,
                    weight: 100,
                }
            ],
            args: [
            ],
            services: [
                {
                    name: "lab-app",
                    port: 3000,
                    path: "/app"
                }
            ],
            healthcheck:
            {
                path: "/healthcheck",
                port: 3000
            },
        },
        slug: "example-lab",
        name: "Example Lab",
        description: "This is an example lab",
        vulnerabilities: ['sql-injection', 'broken-access-control'],
        difficulty: 1,
        languages: ['javascript'],
        technologies: ['nodejs', 'express', 'redis'],
        references: ['https://www.owasp.org/index.php/SQL_Injection', 'https://www.owasp.org/index.php/Broken_Access_Control'],
        authors: ['AppSecLabs'],
        rating: 4,
        estimated_time: 10,
        status: "active",
        created_at: new Date(),
        updated_at: new Date()
    },
    {
        labSpec: {
            image: "ghcr.io/vitor-mauricio/lab-runner-copy:dev",
            ports: [3000],
            code_config: {
                git_url: "github.com/appsec-digital/labs.git",
                git_branch: "main",
                git_path: "copy-n-paste",
            },
            env: {
                "DOCKER_HOST": "tcp://localhost:2375",
                "APP_PORT": "3000",
                "APP_URL": "$(LAB_BASE_URL)/$(NAMESPACE)/proxy/$(APP_PORT)"
            },
            args: [
            ],
            services: [
                {
                    name: "lab-app",
                    port: 3000,
                    path: "/app"
                }
            ],
            evaluations: [
                {
                    name: "sast",
                    order: 1,
                    weight: 100,
                }
            ],
            healthcheck:
            {
                path: "/healthcheck",
                port: 3000
            },

        },
        slug: "copy-n-paste-lab",
        name: "Copy-n-Paste Lab",
        description: "CopyNPaste is a Golang web application that uses an API and a simple front end to simulate a login page. It has both /register and /login routes that, by communicating with a MySQL database, enable users to register and enter into a generic system.",
        vulnerabilities: ['sql-injection', 'injection-flaws'],
        difficulty: 2,
        languages: ['golang', 'html', 'javascript'],
        technologies: ['golang', 'echo', 'mariadb', 'templates'],
        references: ['https://www.owasp.org/index.php/SQL_Injection', 'https://github.com/globocom/secDevLabs/tree/master/owasp-top10-2021-apps/a3/copy-n-paste'],
        authors: ['Globo'],
        rating: 4,
        estimated_time: 30,
        status: "active",
        created_at: new Date(),
        updated_at: new Date()
    }
]);

db.evaluations.insertMany([
    {
        name: "SAST",
        slug: "sast",
        description: "Static Application Security Testing",
        evaluation_spec: {
            containers: [
                {
                    name: "sast",
                    image: "ghcr.io/appsec-digital/eval-sast:dev",
                    env: [],
                    commands: ["/bin/sh", "-c"],
                    args: ["semgrep scan /workspace --json --output /results/result.json >/dev/null 2>&1 && cat /results/result.json | jq '.' || echo 'Semgrep scan failed'"],
                    volumes: [
                    ]
                }
            ],
            init_container: {
            },
            volumes: [
            ]
        },
        created_at: new Date(),
        updated_at: new Date()
    },
    {
    name: "Exploit",
    slug: "exploit",
    description: "Exploit",
    evaluation_spec: {
        containers: [
            {
                name: "exploit",
                image: "ghcr.io/appsec-digital/eval-exploit:dev",
                env: [],
                commands: [
                ],
                args: [
                    "python", "/app/evaluate.py", "/workspace","/template", "/results/result.json"
                ],
                volumes: [
                    {
                        name: "template-volume",
                        mountPath: "/template"
                    }
                ],
                security_context:{
                    privileged: true
                }

            }
        ],
        init_container: {
            name: "clone-tamplates",
            image: "alpine/git",
            commands: [
               "sh", "-c", "set -e; echo 'Cloning lab...'; git clone --depth 1 --branch $GIT_BRANCH https://$GIT_USERNAME:$GIT_PASSWORD@$GIT_URL /tmp/lab; echo 'Copying exploit templates...'; cp -r /tmp/lab/$GIT_PATH/* /template; echo 'Chowning lab...'; chown -R 1000:1000 /template"
            ],
            volumes: [
                {
                    name: "template-volume",
                    mountPath: "/template"
                }
            ],
        },
        volumes: [
            {
                name: "template-volume",
                mountPath: "/template"
            }
        ]
    },
    created_at: new Date(),
    updated_at: new Date()
    }
]);

print('Inicialização do MongoDB concluída com sucesso!'); 