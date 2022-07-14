1. build_chain
	SHORT_SHA
	DOCKER_REPO & IMAGE_NAME

2. deploy_chain
	SHORT_SHA
	HOST
	DOCKER_REPO & IMAGE_NAME
	deploy.sh path

	OUTPUT:
	1. List of env vars -> file/std

3. build_bmr
	SHORT_SHA
	DOCKER_REPO & IMAGE_NAME

4. deploy_contracts
	SHORT_SHA
	OUTPUT of (2) as input

	OUTPUT: (env file combined)
	1. bmr_config(json file)
	2. e2e_config(json file)

5. deploy_bmr
	SHORT_SHA
	DOCKER_REPO & IMAGE_NAME
	OUTPUT of (4)(1) bmr_config

6. build_e2e
	SHORT_SHA
	DOCKER_REPO & IMAGE_NAME

7. deploy_e2e
	OUPPUT of (4)(2) e2e_config
