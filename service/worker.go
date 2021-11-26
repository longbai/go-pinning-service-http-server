package service

import (
	"context"
	"fmt"
	"io"
	"log"

	cid "github.com/ipfs/go-cid"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/ipfs/go-ipfs-http-client"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/longbai/go-pinning-service-http-server/model"
	ma "github.com/multiformats/go-multiaddr"
)

//def ipfs_add
//begin
//update_columns(status: 'pinning') unless status == 'pinned'
//origins.each do |origin|
//ipfs_client.swarm_connect(origin)
//end
//ipfs_client.pin_add(cid)
//update_columns(status: 'pinned')
//rescue Ipfs::Commands::Error => e
//puts e
//# TODO record the exception somewhere
//update_columns(status: 'failed')
//end
//end
//
//def ipfs_remove_async
//IpfsRemoveWorker.perform_async(id)
//end
//
//def ipfs_remove
//# TODO only unpin cid if this is the only pin with that CID
//begin
//ipfs_client.pin_rm(cid)
//update_columns(status: 'removed')
//rescue Ipfs::Commands::Error => e
//raise unless JSON.parse(e.message)['Code'] == 0
//end
//end
type ipfsConfig struct {
	ip string
	port string
}

var ipfsCfg ipfsConfig

func IpfsInit(ip, port string) {
	ipfsCfg.ip = ip
	ipfsCfg.port = port
}

func newIpfsClient(ip, port string)(*httpapi.HttpApi, error) {
	a := fmt.Sprintf("/ip4/%s/tcp/%s", ip, port)
	addr, err := ma.NewMultiaddr(a)
	if err != nil {
		log.Println("ipfs client error", err)
		return nil, err
	}
 	return httpapi.NewApi(addr)
}

func ipfsPinAdd(ctx context.Context, pin *model.PinStatus) error{
	c, err := newIpfsClient(ipfsCfg.ip, ipfsCfg.port)
	if err != nil {
		log.Println("ipfs add client error", err)
		return err
	}
	id, err := cid.Decode(pin.Pin.Cid)
	if err != nil {
		log.Println("ipfs add cid error", err)
		return err
	}
	p := path.IpfsPath(id)
	return c.Pin().Add(ctx, p)
}

func ipfsPinRemove(ctx context.Context, pin *model.PinStatus) error{
	c, err := newIpfsClient(ipfsCfg.ip, ipfsCfg.port)
	if err != nil {
		log.Println("ipfs rem client error", err)
		return err
	}
	id, err := cid.Decode(pin.Pin.Cid)
	if err != nil {
		log.Println("ipfs add cid error", err)
		return err
	}
	p := path.IpfsPath(id)
	return c.Pin().Rm(ctx, p)
}

func ipfsList(ctx context.Context) ([]string, error) {
	c, err := newIpfsClient(ipfsCfg.ip, ipfsCfg.port)
	if err != nil {
		log.Println("ipfs ls client error", err)
		return nil, err
	}
	p, err := c.Pin().Ls(ctx)
	if err != nil {
		log.Println("ipfs ls error", err)
		return nil, err
	}
	var ls []string
	for v := range p {
		ls = append(ls, v.Path().Cid().String())
	}
	return ls, nil
}

func ipfsPut(ctx context.Context, reader io.Reader) (string, error) {
	sh := shell.NewShell(fmt.Sprintf("%s:%s", ipfsCfg.ip, ipfsCfg.port))
	fileHash, err := sh.Add(reader)
	return fileHash, err
	//c, err := newIpfsClient(ipfsCfg.ip, ipfsCfg.port)
	//if err != nil {
	//	log.Println("ipfs ls client error", err)
	//	return "", err
	//}
	//obj, err := c.Object().Put(ctx, reader)
	//if err != nil {
	//	log.Println("ipfs put error", err)
	//	return "", err
	//}
	//return obj.Cid().String(), nil
}
