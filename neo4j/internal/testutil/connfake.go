/*
 * Copyright (c) "Neo4j"
 * Neo4j Sweden AB [http://neo4j.com]
 *
 * This file is part of Neo4j.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package testutil

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/log"
	"time"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
)

type Next struct {
	Record  *db.Record
	Summary *db.Summary
	Err     error
}

type RecordedTx struct {
	Origin    string
	Mode      db.AccessMode
	Bookmarks []string
	Timeout   time.Duration
	Meta      map[string]interface{}
}

type ConnFake struct {
	Name           string
	Version        string
	Alive          bool
	Birth          time.Time
	Table          *db.RoutingTable
	Err            error
	Id             int
	TxBeginErr     error
	TxBeginHandle  db.TxHandle
	RunErr         error
	RunStream      db.StreamHandle
	RunTxErr       error
	RunTxStream    db.StreamHandle
	Nexts          []Next
	Bookm          string
	TxCommitErr    error
	TxCommitHook   func()
	TxRollbackErr  error
	ResetHook      func()
	ConsumeSum     *db.Summary
	ConsumeErr     error
	ConsumeHook    func()
	RecordedTxs    []RecordedTx // Appended to by Run/TxBegin
	BufferErr      error
	BufferHook     func()
	ForceResetHook func() error
	DatabaseName   string
}

func (c *ConnFake) ServerName() string {
	return c.Name
}

func (c *ConnFake) IsAlive() bool {
	return c.Alive
}

func (c *ConnFake) Reset() {
}

func (c *ConnFake) Close() {
}

func (c *ConnFake) Birthdate() time.Time {
	return c.Birth
}

func (c *ConnFake) Bookmark() string {
	return c.Bookm
}

func (c *ConnFake) ServerVersion() string {
	return "serverVersion"
}

func (c *ConnFake) Buffer(streamHandle db.StreamHandle) error {
	if c.BufferHook != nil {
		c.BufferHook()
	}
	return c.BufferErr
}

func (c *ConnFake) Consume(streamHandle db.StreamHandle) (*db.Summary, error) {
	if c.ConsumeHook != nil {
		c.ConsumeHook()
	}
	return c.ConsumeSum, c.ConsumeErr
}

func (c *ConnFake) GetRoutingTable(context map[string]string, bookmarks []string, database, impersonatedUser string) (*db.RoutingTable, error) {
	if c.Table != nil {
		c.Table.DatabaseName = database
	}
	return c.Table, c.Err
}

func (c *ConnFake) TxBegin(txConfig db.TxConfig) (db.TxHandle, error) {
	c.RecordedTxs = append(c.RecordedTxs, RecordedTx{Origin: "TxBegin", Mode: txConfig.Mode, Bookmarks: txConfig.Bookmarks, Timeout: txConfig.Timeout, Meta: txConfig.Meta})
	return c.TxBeginHandle, c.TxBeginErr
}

func (c *ConnFake) TxRollback(tx db.TxHandle) error {
	return c.TxRollbackErr
}

func (c *ConnFake) TxCommit(tx db.TxHandle) error {
	if c.TxCommitHook != nil {
		c.TxCommitHook()
	}
	return c.TxCommitErr
}

func (c *ConnFake) Run(runCommand db.Command, txConfig db.TxConfig) (db.StreamHandle, error) {

	c.RecordedTxs = append(c.RecordedTxs, RecordedTx{Origin: "Run", Mode: txConfig.Mode, Bookmarks: txConfig.Bookmarks, Timeout: txConfig.Timeout, Meta: txConfig.Meta})
	return c.RunStream, c.RunErr
}

func (c *ConnFake) RunTx(tx db.TxHandle, runCommand db.Command) (db.StreamHandle, error) {
	return c.RunTxStream, c.RunTxErr
}

func (c *ConnFake) Keys(streamHandle db.StreamHandle) ([]string, error) {
	return nil, nil
}

func (c *ConnFake) Next(streamHandle db.StreamHandle) (*db.Record, *db.Summary, error) {
	next := c.Nexts[0]
	if len(c.Nexts) > 1 {
		c.Nexts = c.Nexts[1:]
	}
	return next.Record, next.Summary, next.Err
}

func (c *ConnFake) ForceReset() error {
	if c.ForceResetHook != nil {
		return c.ForceResetHook()
	}
	return nil
}

func (c *ConnFake) SelectDatabase(database string) {
	c.DatabaseName = database
}

func (c *ConnFake) SetBoltLogger(_ log.BoltLogger) {
}
