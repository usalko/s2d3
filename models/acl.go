/**
 * Copyright (C) 2024 Vanya Usalko <ivict@rambler.ru>
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * File: acl.go
 */
package models

const (
	PrivateACL                = "private"
	PublicReadACL             = "public-read"
	PublicReadWriteACL        = "public-read-write"
	AWSExecReadACL            = "aws-exec-read"
	AuthenticatedReadACL      = "authenticated-read"
	BucketOwnerReadACL        = "bucket-owner-read"
	BucketOwnerFullControlACL = "bucket-owner-full-control"
	LogDeliveryWriteACL       = "log-delivery-write"
)
