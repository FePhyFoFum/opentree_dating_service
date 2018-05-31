package inducedates

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var tree Tree
var namedNodeMap map[int]*Node
var ottidNameMap map[string]string
var gbifOttidMap map[int]int
var ncbiNameMap map[string]string

// Run at the start of the server
func init() {
	tfn := "labelled_supertree.dated.tiplen.tre"
	//read a tree file
	f, err := os.Open(tfn)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	scanner := bufio.NewReader(f)
	var rt *Node
	for {
		ln, err := scanner.ReadString('\n')
		if len(ln) > 0 {
			rt = ReadNewickString(ln)
			tree.Instantiate(rt)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}
	}
	fmt.Println(len(tree.Tips))
	// need to get all the named nodes into a dictionary
	namedNodeMap = make(map[int]*Node)
	for _, n := range tree.Pre {
		if len(n.Nam) < 3 {
			continue
		}
		if n.Nam[:3] == "ott" {
			i, err := strconv.Atoi(n.Nam[3:])
			if err == nil {
				namedNodeMap[i] = n
			} else {
				fmt.Println(err)
			}
		}
	}
	fmt.Println("finished tree processing")
	//end the tree tracing information
	//start the taxonomy information
	tfn = "taxonomy.tsv"
	f, err = os.Open(tfn)
	defer f.Close()
	scanner = bufio.NewReader(f)
	ottidNameMap = make(map[string]string)
	gbifOttidMap = make(map[int]int)
	for {
		ln, err := scanner.ReadString('\n')
		if len(ln) > 0 {
			s := strings.Split(ln, "\t|\t")
			ottidNameMap[s[0]] = prettifyName(s[2])
			if strings.Contains(s[4], "gbif") {
				s2 := strings.Split(s[4], ",")
				var gbifid int
				for _, x := range s2 {
					if strings.Contains(x, "gbif") {
						s3 := strings.Split(x, ":")
						nv, err := strconv.Atoi(s3[1])
						if err != nil {
							fmt.Println(err)
							os.Exit(0)
						}
						gbifid = nv
						break
					}
				}
				otdid, err := strconv.Atoi(s[0])
				if err != nil {
					fmt.Println(err)
					os.Exit(0)
				}
				gbifOttidMap[gbifid] = otdid
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}
	}
	fmt.Println("finished tax processing")
	//end the taxonomy information
	//start ncbi taxonomy processing
	tfn = "names.dmp"
	f, err = os.Open(tfn)
	defer f.Close()
	scanner = bufio.NewReader(f)
	ncbiNameMap = make(map[string]string)
	for {
		ln, err := scanner.ReadString('\n')
		if len(ln) > 0 {
			s := strings.Split(ln, "\t|\t")
			if s[3] == "scientific name\t|\n" {
				ncbiNameMap[s[0]] = prettifyName(s[1])
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}
	}
	fmt.Println("finished tax (ncbi) processing")
	//end ncbi taxonomy processing

	fmt.Println("finished init")
}

//GetOttidsFromGbifids gets ott ids from gbif ids
func GetOttidsFromGbifids(ids []int) OttidResults {
	var n OttidResults
	n.Ottids = make([]int, 0)
	n.Unmatched = make([]int, 0)
	for _, i := range ids {
		if _, ok := gbifOttidMap[i]; ok {
			n.Ottids = append(n.Ottids, gbifOttidMap[i])
		} else {
			n.Unmatched = append(n.Unmatched, i)
		}
	}
	return n
}

//GetInducedTree make the induced tree and return the newick string
func GetInducedTree(ids []int) Newick {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rid := r.Float64()
	nds := make([]*Node, 0)
	unmatched := make([]int, 0)
	for _, i := range ids {
		if _, ok := namedNodeMap[i]; ok {
			nds = append(nds, namedNodeMap[i])
		} else {
			unmatched = append(unmatched, i)
		}
	}
	traceTree(nds, rid)
	addLen := make(map[string]float64)
	for _, i := range nds {
		if len(i.Chs) > 0 {
			if len(getMarkedChs(i, rid)) == 0 {
				addLen[i.Nam] = getSubTendLen(i)
			}
		}
	}
	var n Newick
	n.Unmatched = unmatched
	curRt := tree.Rt
	going := true
	for going {
		x := getMarkedChs(curRt, rid)
		if len(x) == 1 {
			curRt = x[0]
		} else {
			going = false
			break
		}
	}
	x := curRt.NewickPaint(true, rid) + ";"
	//handle the tips thjat are not tips
	//this takes time but necessary
	nt := ReadNewickString(x)
	var ntt Tree
	ntt.Instantiate(nt)
	for _, i := range ntt.Tips {
		if _, ok := addLen[i.Nam]; ok {
			i.Len += addLen[i.Nam]
		}
	}
	x = nt.Newick(true) + ";"
	//end the handle
	n.NewString = x
	untraceTree(nds, rid)
	return n
}

func getSubTendLen(nd *Node) float64 {
	x := 0.0
	going := true
	curNd := nd.Chs[0]
	for going {
		x += curNd.Len
		if len(curNd.Chs) > 0 {
			curNd = curNd.Chs[0]
		} else {
			going = false
			break
		}
	}
	return x
}

func getMarkedChs(nd *Node, rid float64) []*Node {
	x := make([]*Node, 0)
	for _, i := range nd.Chs {
		if _, ok := i.MarkedMap[rid]; ok {
			x = append(x, i)
		}
	}
	return x
}

func traceTree(nds []*Node, rid float64) {
	for _, n := range nds {
		n.MarkedMap[rid] = true
		going := true
		cur := n.Par
		for going {
			if _, ok := cur.MarkedMap[rid]; ok {
				break
			}
			cur.MarkedMap[rid] = true
			if cur.Par == nil {
				break
			}
			cur = cur.Par
		}
	}
}

func untraceTree(nds []*Node, rid float64) {
	for _, n := range nds {
		if _, ok := n.MarkedMap[rid]; ok {
			delete(n.MarkedMap, rid)
		}
		going := true
		cur := n.Par
		for going {
			if _, ok := cur.MarkedMap[rid]; !ok {
				break
			}
			delete(cur.MarkedMap, rid)
			if cur.Par == nil {
				break
			}
			cur = cur.Par
		}
	}
}

//GetRenamedTree innewick with ottids and out newick with ottids
func GetRenamedTree(newickin string) Newick {
	var renew Newick
	t := ReadNewickString(newickin)
	var intree Tree
	intree.Instantiate(t)
	for _, i := range intree.Post {
		if len(i.Nam) == 0 {
			continue
		}
		mn := i.Nam
		if i.Nam[:3] == "ott" {
			mn = i.Nam[3:]
		}
		if _, ok := ottidNameMap[mn]; ok {
			i.Nam = ottidNameMap[mn]
		}
	}
	renew.NewString = intree.Rt.Newick(true) + ";"
	return renew
}

//GetRenamedTreeNCBI innewick with ottids and out newick with ncbiids
func GetRenamedTreeNCBI(newickin string) Newick {
	var renew Newick
	t := ReadNewickString(newickin)
	var intree Tree
	intree.Instantiate(t)
	for _, i := range intree.Post {
		if len(i.Nam) == 0 {
			continue
		}
		mn := i.Nam
		if _, ok := ncbiNameMap[mn]; ok {
			i.Nam = ncbiNameMap[mn]
		}
	}
	renew.NewString = intree.Rt.Newick(true) + ";"
	return renew
}

func prettifyName(ins string) string {
	x := ins
	x = strings.Replace(x, "\"", "", -1)
	x = strings.Replace(x, "'", "", -1)
	x = strings.Replace(x, ";", "", -1)
	x = strings.Replace(x, "(", "", -1)
	x = strings.Replace(x, ")", "", -1)
	x = strings.Replace(x, ":", "", -1)
	x = strings.Replace(x, ",", "", -1)
	x = strings.Replace(x, " ", "_", -1)
	x = strings.Replace(x, "[", "", -1)
	x = strings.Replace(x, "]", "", -1)
	x = strings.Replace(x, "&", "", -1)
	return x
}
