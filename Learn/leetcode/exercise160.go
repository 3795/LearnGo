package main

import "fmt"

/**
相交链表
https://leetcode-cn.com/problems/intersection-of-two-linked-lists/
*/

type ListNode struct {
	Val  int
	Next *ListNode
}

func genListNode(arr []int) *ListNode {
	head := &ListNode{arr[0], nil}
	cur := head
	for i := 1; i < len(arr); i++ {
		cur.Next = &ListNode{arr[i], nil}
		cur = cur.Next
	}
	return head
}

func printListNode(head *ListNode) {
	cur := head
	for {
		if cur == nil {
			break
		}
		fmt.Printf("%d -> ", cur.Val)
		cur = cur.Next
	}
}

func getIntersectionNode(headA, headB *ListNode) *ListNode {
	lengthA := getLength(headA)
	lengthB := getLength(headB)
	curA := headA
	curB := headB
	for {
		if lengthA == lengthB {
			break
		}

		if lengthA > lengthB {
			curA = curA.Next
			lengthA--
		} else {
			curB = curB.Next
			lengthB--
		}
	}

	for i := 0; i < lengthA; i++ {
		if curA.Val == curB.Val {
			return curA
		}
		curA = curA.Next
		curB = curB.Next
	}
	return nil
}

func getLength(head *ListNode) int {
	cur := head
	length := 0
	for {
		if cur == nil {
			break
		}
		length++
		cur = cur.Next
	}
	return length
}

func main() {
	arrA := []int{4, 1, 8, 4, 5}
	arrB := []int{5, 0, 1, 8, 4, 5}
	headA := genListNode(arrA)
	headB := genListNode(arrB)
	result := getIntersectionNode(headA, headB)
	fmt.Printf("val: %d", result.Val)
}
