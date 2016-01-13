namespace go apm

struct GoBin {
  1: string SourcePath,
  2: string Name,
  3: bool KeepAlive = false,
  4: list<string> Args,
}

struct ProcStatus {
  1: string Status,
  2: i32 Restarts,
}

struct Proc {
  1: string Name,
  2: string Cmd,
  3: list<string> Args,
  4: string Path,
  5: string Pidfile,
  6: string Outfile,
  7: string Errfile,
  8: bool KeepAlive,
  9: i32 Pid,
  10: ProcStatus Status,
}

service Apm {

  void ping(),

  void save(),

  void resurrect(),

  string gobin(1:GoBin goBin),

  void startProc(1:string procName),

  void stopProc(1:string procName),

  void restartProc(1:string procName),

  void deleteProc(1:string procName),

  list<Proc> monit()

}
  