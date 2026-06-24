
typedef enum { EnvAdd, EnvRemove, EnvUpdate } EnvAction;

// Args struct
typedef struct {
  int help;
  int addEnv;
  const char *envName;
  const char *envValue;
  EnvAction action;
} Args;

// Parse the arguments
int parseArgs(Args *args, int argc, char *argv[]);
