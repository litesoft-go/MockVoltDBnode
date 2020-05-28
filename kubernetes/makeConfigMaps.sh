#!/bin/bash -E

err_report() {
    echo "Exit on line $(caller)" >&2
    exit 111
}

trap err_report ERR

build_configmap_from_files() {
    local MAPNAME=$1
    shift 1
    # Takes a comma-separated list of file paths, which may include Unix-style wildcards (*?[ch-ranges])
    # returns a comma separated list of all matching files, following links
    # nb. find may fail, inform user
    if [[ -n $@ ]]; then
        local FILES=$(find ${@//,/ } -type f -print)
        local RC=$?
        if [[ $RC -ne 0 ]]; then
            echo "ERROR: create configmap '$MAPNAME' failed due to error(s)" >&2
            exit 111
        fi
    fi
    local ORDERFILE=$(mktemp)
    local CONFIG_MAP_ARGS=""
    if [[ -n $FILES ]]; then
        [[ -n $FILES ]] && echo "Files found for configmap/$MAPNAME:"
        local FNS=""
        for f in $FILES
        do
            echo -e "\t$f"
            CONFIG_MAP_ARGS+="--from-file=${f} "
            FNS+="$(basename ${f}),"
        done
        echo "$FNS" | sed -e 's/,$//' > $ORDERFILE
        CONFIG_MAP_ARGS+="--from-file=.loadorder=$ORDERFILE"
    fi
    echo "Creating configmap '$MAPNAME'.."
    # if there are no files, we create an empty configmap
    kubectl_ifexists delete configmap "$MAPNAME"
    $KUBECTL create configmap "$MAPNAME" $CONFIG_MAP_ARGS
    rm -f $ORDERFILE
}


kubectl_ifexists() {
    RC=0
    if $KUBECTL get "$2"/"$3" &>/dev/null; then
        $KUBECTL $@
        RC=$?
    fi
    return $RC
}


make_configmaps() {

    # make classes configmap
    # packages voltdb procedures/udf classes jar files
    MAPNAME=${CLUSTER_NAME}-classes
    build_configmap_from_files $MAPNAME "${CLASSES_FILES}"

    # make schema configmap
    # packages voltdb ddl files
    MAPNAME=${CLUSTER_NAME}-schema
    kubectl_ifexists delete configmap "$MAPNAME"
    build_configmap_from_files $MAPNAME "${SCHEMA_FILES}"

    # make init configmap
    # packages: deployment.xml, license.xml, log4j properties
    MAPNAME=${CLUSTER_NAME}-init
    kubectl_ifexists delete configmap "$MAPNAME"
    CONFIG_MAP_ARGS=""
    [[ -n ${DEPLOYMENT_FILE} ]] && CONFIG_MAP_ARGS+=" --from-file=deployment=${DEPLOYMENT_FILE}"
    [[ -n ${LICENSE_FILE} ]] && CONFIG_MAP_ARGS+=" --from-file=license=${LICENSE_FILE}"
    [[ -n ${LOG4J_CONFIG_PATH} ]] && CONFIG_MAP_ARGS+=" --from-file=log4jcfg=${LOG4J_CONFIG_PATH}"

    $KUBECTL create configmap $MAPNAME $CONFIG_MAP_ARGS

    # make runtime environment settings configmap
    # packages voltdb runtime settings NODECOUNT, VOLTDB_OPTS, etc.
    MAPNAME=${CLUSTER_NAME}-env
    kubectl_ifexists delete configmap "$MAPNAME"

    PROP_FILE="$(mktemp)"
    CONFIG_MAP_ARGS="--from-env-file=${PROP_FILE}"
    [[ -n ${VOLTDB_INIT_ARGS} ]] && echo "VOLTDB_INIT_ARGS=${VOLTDB_INIT_ARGS}"     >> ${PROP_FILE}
    grep VOLTDB_START_ARGS      ${CFG_FILE} >> ${PROP_FILE}
    egrep NODECOUNT             ${CFG_FILE} >> ${PROP_FILE}
    [[ -n ${LOG4J_CONFIG_PATH}  ]] && echo "LOG4J_CONFIG_PATH=\"/etc/voltdb/log4jcfg\"" >> ${PROP_FILE}
    [[ -n ${VOLTDB_OPTS} ]] && echo "VOLTDB_OPTS=${VOLTDB_OPTS}"                    >> ${PROP_FILE}
    grep VOLTDB_HEAPMAX         ${CFG_FILE} >> ${PROP_FILE}

    $KUBECTL create configmap $MAPNAME --from-env-file=${PROP_FILE}

    rm -f ${PROP_FILE}
}

# MAIN

# source the template settings
if [[ ! -e "$1" ]]; then
    echo "ERROR config file not specified; customize config-template.cfg with your settings and database assets"
    exit 1
fi

# source our config file parameters
CFG_FILE=$1
source ${CFG_FILE}

shift 1

# Hook to support alternative kubectl commands if necessary
: ${KUBECTL:=kubectl}

# use config file name as default cluster name
: ${CLUSTER_NAME:=$(basename $CFG_FILE ".cfg")}

make_configmaps
