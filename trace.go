package inducedates

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

var tree Tree
var namedNodeMap map[int]*Node

// Give us some seed data
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
	fmt.Println("finished init")
}

//GetInducedTree make the induced tree and return the newick string
func GetInducedTree(ids []int) Newick {
	nds := make([]*Node, 0)
	unmatched := make([]int, 0)
	for _, i := range ids {
		if _, ok := namedNodeMap[i]; ok {
			nds = append(nds, namedNodeMap[i])
		} else {
			unmatched = append(unmatched, i)
		}
	}
	traceTree(nds)
	addLen := make(map[string]float64)
	for _, i := range nds {
		if len(i.Chs) > 0 {
			if len(getMarkedChs(i)) == 0 {
				addLen[i.Nam] = getSubTendLen(i)
			}
		}
	}
	var n Newick
	n.Unmatched = unmatched
	curRt := tree.Rt
	going := true
	for going {
		x := getMarkedChs(curRt)
		if len(x) == 1 {
			curRt = x[0]
		} else {
			going = false
			break
		}
	}
	x := curRt.NewickPaint(true) + ";"
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
	untraceTree(nds)
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

func getMarkedChs(nd *Node) []*Node {
	x := make([]*Node, 0)
	for _, i := range nd.Chs {
		if i.Marked {
			x = append(x, i)
		}
	}
	return x
}

func traceTree(nds []*Node) {
	for _, n := range nds {
		n.Marked = true
		going := true
		cur := n.Par
		for going {
			if cur.Marked {
				break
			}
			cur.Marked = true
			if cur.Par == nil {
				break
			}
			cur = cur.Par
		}
	}
}

func untraceTree(nds []*Node) {
	for _, n := range nds {
		n.Marked = false
		going := true
		cur := n.Par
		for going {
			if cur.Marked == false {
				break
			}
			cur.Marked = false
			if cur.Par == nil {
				break
			}
			cur = cur.Par
		}
	}
}
