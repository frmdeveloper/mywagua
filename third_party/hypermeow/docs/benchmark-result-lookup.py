"""Bounded isolated read-recovery benchmark. Run on Linux with Docker and cgroup v2.
Usage: python3 benchmark-result-lookup.py VERIFIED_LOCAL_IMAGE
Only a new unlinked, unpaired fixture is modified. No production allocation is queried.
"""
import concurrent.futures, json, math, pathlib, secrets, subprocess, sys, time, urllib.request, urllib.error, uuid

def run(*args):
    return subprocess.check_output(args, text=True).strip()

def main():
    image = sys.argv[1]
    ident = uuid.uuid4().hex
    name = "wpp-bench-" + ident
    label = "wppconnect.benchmark=" + ident
    metadata = json.loads(run("docker","image","inspect",image))[0]
    created = False
    try:
        run("docker","run","-d","--name",name,"--label",label,"--memory","256m","--cpus","0.75","--pids-limit","96","-p","127.0.0.1::21465","-e","AUTO_RESTORE=false","-e","SECRET_KEY="+secrets.token_hex(24),metadata["Id"])
        created = True
        port = run("docker","port",name,"21465/tcp").split(":")[-1]
        base = "http://127.0.0.1:"+port
        def call(path, token=None, key=None, body=None):
            headers = {"Content-Type":"application/json"}
            if token: headers["Authorization"]="Bearer "+token
            if key: headers["Idempotency-Key"]=key
            try:
                with urllib.request.urlopen(urllib.request.Request(base+path,data=body,headers=headers),timeout=12) as r:
                    return r.status,r.read()
            except urllib.error.HTTPError as e:
                return e.code,e.read()
        for _ in range(40):
            try:
                if call("/readyz")[0]==200: break
            except OSError: pass
            time.sleep(.25)
        else: raise RuntimeError("fixture was not ready")
        token=json.loads(run("docker","exec",name,"/app/token","-session","bench"))["token"]
        other=json.loads(run("docker","exec",name,"/app/token","-session","other"))["token"]
        seeds=[]
        for i in range(4):
            key="sample-key-"+str(i)
            code,body=call("/api/bench/newsletter-mark-viewed",token,key,b"{}")
            assert code==404
            code,raw=call("/api/bench/newsletter-write-result",token,key)
            assert code==200 and json.loads(raw)["httpStatus"]==404
            assert call("/api/other/newsletter-write-result",other,key)[0]==404
            seeds.append((key,raw))
        inspect=json.loads(run("docker","inspect",name))[0]
        pid=inspect["State"]["Pid"]
        cgroup=next(line.split("::",1)[1] for line in pathlib.Path(f"/proc/{pid}/cgroup").read_text().splitlines() if line.startswith("0::"))
        cg=pathlib.Path("/sys/fs/cgroup")/cgroup.lstrip("/")
        def cpu():
            return int(next(line.split()[1] for line in (cg/"cpu.stat").read_text().splitlines() if line.startswith("usage_usec ")))
        def lookup(i):
            key,saved=seeds[i%len(seeds)]
            start=time.monotonic();code,body=call("/api/bench/newsletter-write-result",token,key);dt=time.monotonic()-start
            assert code==200 and body==saved
            return dt*1000
        profiles=[]
        for workers,n in [(1,24),(4,48)]:
            c0=cpu(); start=time.monotonic()
            with concurrent.futures.ThreadPoolExecutor(max_workers=workers) as pool: times=sorted(pool.map(lookup,range(n)))
            wall=time.monotonic()-start
            def percentile(p):return round(times[math.ceil(p*n)-1],2)
            profiles.append({"workers":workers,"requests":n,"byteIdenticalResponses":n,"p50Ms":percentile(.5),"p95Ms":percentile(.95),"p99Ms":percentile(.99),"requestsPerSecond":round(n/wall,2),"cpuSeconds":round((cpu()-c0)/1e6,3),"peakCgroupMemoryBytes":int((cg/"memory.peak").read_text())})
        run("docker","restart",name)
        base = "http://127.0.0.1:"+run("docker","port",name,"21465/tcp").split(":")[-1]
        for _ in range(40):
            try:
                if call("/readyz")[0]==200:break
            except OSError:pass
            time.sleep(.25)
        else: raise RuntimeError("fixture did not become ready after restart")
        for i in range(4):lookup(i)
        result={"fixtureImage":metadata["Id"],"source":metadata["Config"]["Labels"]["org.opencontainers.image.revision"],"cpuLimit":0.75,"memoryLimitMiB":256,"workload":"authenticated lookup of four recorded disconnected-channel results","profiles":profiles,"recordsRecoveredAfterRestart":4,"crossSessionLookupsReturned404":4,"realWhatsAppMessages":0}
    finally:
        if created:
            item=json.loads(run("docker","inspect",name))[0]
            assert item["Config"]["Labels"].get("wppconnect.benchmark")==ident
            run("docker","rm","-fv",name)
    result["fixtureRemoved"]=True
    print(json.dumps(result,indent=2))
if __name__=="__main__":main()
