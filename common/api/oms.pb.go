// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        v5.29.3
// source: api/oms.proto

package api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Order struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ID            string                 `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	CustomerID    string                 `protobuf:"bytes,2,opt,name=customerID,proto3" json:"customerID,omitempty"`
	Status        string                 `protobuf:"bytes,3,opt,name=Status,proto3" json:"Status,omitempty"`
	Items         []*Item                `protobuf:"bytes,4,rep,name=Items,proto3" json:"Items,omitempty"`
	PaymentLink   string                 `protobuf:"bytes,5,opt,name=PaymentLink,proto3" json:"PaymentLink,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Order) Reset() {
	*x = Order{}
	mi := &file_api_oms_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Order) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Order) ProtoMessage() {}

func (x *Order) ProtoReflect() protoreflect.Message {
	mi := &file_api_oms_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Order.ProtoReflect.Descriptor instead.
func (*Order) Descriptor() ([]byte, []int) {
	return file_api_oms_proto_rawDescGZIP(), []int{0}
}

func (x *Order) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Order) GetCustomerID() string {
	if x != nil {
		return x.CustomerID
	}
	return ""
}

func (x *Order) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *Order) GetItems() []*Item {
	if x != nil {
		return x.Items
	}
	return nil
}

func (x *Order) GetPaymentLink() string {
	if x != nil {
		return x.PaymentLink
	}
	return ""
}

type Product struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ID            string                 `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	TaxCode       string                 `protobuf:"bytes,3,opt,name=TaxCode,proto3" json:"TaxCode,omitempty"`
	Description   string                 `protobuf:"bytes,4,opt,name=Description,proto3" json:"Description,omitempty"`
	Active        bool                   `protobuf:"varint,5,opt,name=Active,proto3" json:"Active,omitempty"`
	Price         float32                `protobuf:"fixed32,6,opt,name=Price,proto3" json:"Price,omitempty"`
	Quantity      int64                  `protobuf:"varint,7,opt,name=Quantity,proto3" json:"Quantity,omitempty"`
	Currency      string                 `protobuf:"bytes,8,opt,name=Currency,proto3" json:"Currency,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Product) Reset() {
	*x = Product{}
	mi := &file_api_oms_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Product) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Product) ProtoMessage() {}

func (x *Product) ProtoReflect() protoreflect.Message {
	mi := &file_api_oms_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Product.ProtoReflect.Descriptor instead.
func (*Product) Descriptor() ([]byte, []int) {
	return file_api_oms_proto_rawDescGZIP(), []int{1}
}

func (x *Product) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Product) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Product) GetTaxCode() string {
	if x != nil {
		return x.TaxCode
	}
	return ""
}

func (x *Product) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Product) GetActive() bool {
	if x != nil {
		return x.Active
	}
	return false
}

func (x *Product) GetPrice() float32 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *Product) GetQuantity() int64 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

func (x *Product) GetCurrency() string {
	if x != nil {
		return x.Currency
	}
	return ""
}

type CheckIfItemsInStockRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Items         []*ItemsWithQuantity   `protobuf:"bytes,1,rep,name=Items,proto3" json:"Items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CheckIfItemsInStockRequest) Reset() {
	*x = CheckIfItemsInStockRequest{}
	mi := &file_api_oms_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CheckIfItemsInStockRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckIfItemsInStockRequest) ProtoMessage() {}

func (x *CheckIfItemsInStockRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_oms_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckIfItemsInStockRequest.ProtoReflect.Descriptor instead.
func (*CheckIfItemsInStockRequest) Descriptor() ([]byte, []int) {
	return file_api_oms_proto_rawDescGZIP(), []int{2}
}

func (x *CheckIfItemsInStockRequest) GetItems() []*ItemsWithQuantity {
	if x != nil {
		return x.Items
	}
	return nil
}

type CheckIfItemsInStockResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	InStock       bool                   `protobuf:"varint,1,opt,name=InStock,proto3" json:"InStock,omitempty"`
	Items         []*Item                `protobuf:"bytes,2,rep,name=Items,proto3" json:"Items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CheckIfItemsInStockResponse) Reset() {
	*x = CheckIfItemsInStockResponse{}
	mi := &file_api_oms_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CheckIfItemsInStockResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckIfItemsInStockResponse) ProtoMessage() {}

func (x *CheckIfItemsInStockResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_oms_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckIfItemsInStockResponse.ProtoReflect.Descriptor instead.
func (*CheckIfItemsInStockResponse) Descriptor() ([]byte, []int) {
	return file_api_oms_proto_rawDescGZIP(), []int{3}
}

func (x *CheckIfItemsInStockResponse) GetInStock() bool {
	if x != nil {
		return x.InStock
	}
	return false
}

func (x *CheckIfItemsInStockResponse) GetItems() []*Item {
	if x != nil {
		return x.Items
	}
	return nil
}

type GetItemsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ItemIDs       []string               `protobuf:"bytes,1,rep,name=ItemIDs,proto3" json:"ItemIDs,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetItemsRequest) Reset() {
	*x = GetItemsRequest{}
	mi := &file_api_oms_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetItemsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetItemsRequest) ProtoMessage() {}

func (x *GetItemsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_oms_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetItemsRequest.ProtoReflect.Descriptor instead.
func (*GetItemsRequest) Descriptor() ([]byte, []int) {
	return file_api_oms_proto_rawDescGZIP(), []int{4}
}

func (x *GetItemsRequest) GetItemIDs() []string {
	if x != nil {
		return x.ItemIDs
	}
	return nil
}

type GetItemsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Items         []*Item                `protobuf:"bytes,1,rep,name=Items,proto3" json:"Items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetItemsResponse) Reset() {
	*x = GetItemsResponse{}
	mi := &file_api_oms_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetItemsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetItemsResponse) ProtoMessage() {}

func (x *GetItemsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_oms_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetItemsResponse.ProtoReflect.Descriptor instead.
func (*GetItemsResponse) Descriptor() ([]byte, []int) {
	return file_api_oms_proto_rawDescGZIP(), []int{5}
}

func (x *GetItemsResponse) GetItems() []*Item {
	if x != nil {
		return x.Items
	}
	return nil
}

type GetOrderRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	OrderID       string                 `protobuf:"bytes,1,opt,name=OrderID,proto3" json:"OrderID,omitempty"`
	CustomerID    string                 `protobuf:"bytes,2,opt,name=CustomerID,proto3" json:"CustomerID,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetOrderRequest) Reset() {
	*x = GetOrderRequest{}
	mi := &file_api_oms_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetOrderRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOrderRequest) ProtoMessage() {}

func (x *GetOrderRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_oms_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOrderRequest.ProtoReflect.Descriptor instead.
func (*GetOrderRequest) Descriptor() ([]byte, []int) {
	return file_api_oms_proto_rawDescGZIP(), []int{6}
}

func (x *GetOrderRequest) GetOrderID() string {
	if x != nil {
		return x.OrderID
	}
	return ""
}

func (x *GetOrderRequest) GetCustomerID() string {
	if x != nil {
		return x.CustomerID
	}
	return ""
}

type Item struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ID            string                 `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	Quantity      int32                  `protobuf:"varint,3,opt,name=Quantity,proto3" json:"Quantity,omitempty"`
	PriceID       string                 `protobuf:"bytes,4,opt,name=PriceID,proto3" json:"PriceID,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Item) Reset() {
	*x = Item{}
	mi := &file_api_oms_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Item) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Item) ProtoMessage() {}

func (x *Item) ProtoReflect() protoreflect.Message {
	mi := &file_api_oms_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Item.ProtoReflect.Descriptor instead.
func (*Item) Descriptor() ([]byte, []int) {
	return file_api_oms_proto_rawDescGZIP(), []int{7}
}

func (x *Item) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Item) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Item) GetQuantity() int32 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

func (x *Item) GetPriceID() string {
	if x != nil {
		return x.PriceID
	}
	return ""
}

type ItemsWithQuantity struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ID            string                 `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Quantity      int32                  `protobuf:"varint,2,opt,name=Quantity,proto3" json:"Quantity,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ItemsWithQuantity) Reset() {
	*x = ItemsWithQuantity{}
	mi := &file_api_oms_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ItemsWithQuantity) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ItemsWithQuantity) ProtoMessage() {}

func (x *ItemsWithQuantity) ProtoReflect() protoreflect.Message {
	mi := &file_api_oms_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ItemsWithQuantity.ProtoReflect.Descriptor instead.
func (*ItemsWithQuantity) Descriptor() ([]byte, []int) {
	return file_api_oms_proto_rawDescGZIP(), []int{8}
}

func (x *ItemsWithQuantity) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *ItemsWithQuantity) GetQuantity() int32 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

type CreateOrderRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CustomerID    string                 `protobuf:"bytes,1,opt,name=customerID,proto3" json:"customerID,omitempty"`
	Items         []*ItemsWithQuantity   `protobuf:"bytes,2,rep,name=Items,proto3" json:"Items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateOrderRequest) Reset() {
	*x = CreateOrderRequest{}
	mi := &file_api_oms_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateOrderRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateOrderRequest) ProtoMessage() {}

func (x *CreateOrderRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_oms_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateOrderRequest.ProtoReflect.Descriptor instead.
func (*CreateOrderRequest) Descriptor() ([]byte, []int) {
	return file_api_oms_proto_rawDescGZIP(), []int{9}
}

func (x *CreateOrderRequest) GetCustomerID() string {
	if x != nil {
		return x.CustomerID
	}
	return ""
}

func (x *CreateOrderRequest) GetItems() []*ItemsWithQuantity {
	if x != nil {
		return x.Items
	}
	return nil
}

var File_api_oms_proto protoreflect.FileDescriptor

var file_api_oms_proto_rawDesc = string([]byte{
	0x0a, 0x0d, 0x61, 0x70, 0x69, 0x2f, 0x6f, 0x6d, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x03, 0x61, 0x70, 0x69, 0x22, 0x92, 0x01, 0x0a, 0x05, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x0e,
	0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x12, 0x1e,
	0x0a, 0x0a, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x49, 0x44, 0x12, 0x16,
	0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1f, 0x0a, 0x05, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x18,
	0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x49, 0x74, 0x65, 0x6d,
	0x52, 0x05, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x50, 0x61, 0x79, 0x6d, 0x65,
	0x6e, 0x74, 0x4c, 0x69, 0x6e, 0x6b, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x50, 0x61,
	0x79, 0x6d, 0x65, 0x6e, 0x74, 0x4c, 0x69, 0x6e, 0x6b, 0x22, 0xcf, 0x01, 0x0a, 0x07, 0x50, 0x72,
	0x6f, 0x64, 0x75, 0x63, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x54, 0x61, 0x78,
	0x43, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x54, 0x61, 0x78, 0x43,
	0x6f, 0x64, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x50, 0x72, 0x69, 0x63, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x02, 0x52, 0x05, 0x50, 0x72,
	0x69, 0x63, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x51, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x51, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12,
	0x1a, 0x0a, 0x08, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x22, 0x4a, 0x0a, 0x1a, 0x43,
	0x68, 0x65, 0x63, 0x6b, 0x49, 0x66, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x49, 0x6e, 0x53, 0x74, 0x6f,
	0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2c, 0x0a, 0x05, 0x49, 0x74, 0x65,
	0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x49,
	0x74, 0x65, 0x6d, 0x73, 0x57, 0x69, 0x74, 0x68, 0x51, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x52, 0x05, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x22, 0x58, 0x0a, 0x1b, 0x43, 0x68, 0x65, 0x63, 0x6b,
	0x49, 0x66, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x49, 0x6e, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x49, 0x6e, 0x53, 0x74, 0x6f, 0x63,
	0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x49, 0x6e, 0x53, 0x74, 0x6f, 0x63, 0x6b,
	0x12, 0x1f, 0x0a, 0x05, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x09, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x05, 0x49, 0x74, 0x65, 0x6d,
	0x73, 0x22, 0x2b, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x49, 0x74, 0x65, 0x6d, 0x49, 0x44, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x49, 0x74, 0x65, 0x6d, 0x49, 0x44, 0x73, 0x22, 0x33,
	0x0a, 0x10, 0x47, 0x65, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x1f, 0x0a, 0x05, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x09, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x05, 0x49, 0x74,
	0x65, 0x6d, 0x73, 0x22, 0x4b, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x49,
	0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x44,
	0x12, 0x1e, 0x0a, 0x0a, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x49, 0x44, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x49, 0x44,
	0x22, 0x60, 0x0a, 0x04, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08,
	0x51, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08,
	0x51, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x50, 0x72, 0x69, 0x63,
	0x65, 0x49, 0x44, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x50, 0x72, 0x69, 0x63, 0x65,
	0x49, 0x44, 0x22, 0x3f, 0x0a, 0x11, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x57, 0x69, 0x74, 0x68, 0x51,
	0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x12, 0x1a, 0x0a, 0x08, 0x51, 0x75, 0x61, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x51, 0x75, 0x61, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x22, 0x62, 0x0a, 0x12, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4f, 0x72, 0x64,
	0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x75, 0x73,
	0x74, 0x6f, 0x6d, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63,
	0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x49, 0x44, 0x12, 0x2c, 0x0a, 0x05, 0x49, 0x74, 0x65,
	0x6d, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x49,
	0x74, 0x65, 0x6d, 0x73, 0x57, 0x69, 0x74, 0x68, 0x51, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x52, 0x05, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x32, 0x97, 0x01, 0x0a, 0x0c, 0x4f, 0x72, 0x64, 0x65,
	0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x32, 0x0a, 0x0b, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x17, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x0a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x2c, 0x0a, 0x08,
	0x47, 0x65, 0x74, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x14, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x47,
	0x65, 0x74, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0a,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x25, 0x0a, 0x0b, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x0a, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x4f, 0x72, 0x64, 0x65, 0x72, 0x1a, 0x0a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4f, 0x72, 0x64, 0x65,
	0x72, 0x32, 0xa1, 0x01, 0x0a, 0x0c, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x58, 0x0a, 0x13, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x49, 0x66, 0x49, 0x74, 0x65,
	0x6d, 0x73, 0x49, 0x6e, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x12, 0x1f, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x43, 0x68, 0x65, 0x63, 0x6b, 0x49, 0x66, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x49, 0x6e, 0x53, 0x74,
	0x6f, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x49, 0x66, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x49, 0x6e, 0x53,
	0x74, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x37, 0x0a, 0x08,
	0x47, 0x65, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x14, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x47,
	0x65, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x47, 0x65, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x1f, 0x5a, 0x1d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x6a, 0x75, 0x78, 0x75, 0x65, 0x39, 0x37, 0x2f, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_api_oms_proto_rawDescOnce sync.Once
	file_api_oms_proto_rawDescData []byte
)

func file_api_oms_proto_rawDescGZIP() []byte {
	file_api_oms_proto_rawDescOnce.Do(func() {
		file_api_oms_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_oms_proto_rawDesc), len(file_api_oms_proto_rawDesc)))
	})
	return file_api_oms_proto_rawDescData
}

var file_api_oms_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_api_oms_proto_goTypes = []any{
	(*Order)(nil),                       // 0: api.Order
	(*Product)(nil),                     // 1: api.Product
	(*CheckIfItemsInStockRequest)(nil),  // 2: api.CheckIfItemsInStockRequest
	(*CheckIfItemsInStockResponse)(nil), // 3: api.CheckIfItemsInStockResponse
	(*GetItemsRequest)(nil),             // 4: api.GetItemsRequest
	(*GetItemsResponse)(nil),            // 5: api.GetItemsResponse
	(*GetOrderRequest)(nil),             // 6: api.GetOrderRequest
	(*Item)(nil),                        // 7: api.Item
	(*ItemsWithQuantity)(nil),           // 8: api.ItemsWithQuantity
	(*CreateOrderRequest)(nil),          // 9: api.CreateOrderRequest
}
var file_api_oms_proto_depIdxs = []int32{
	7,  // 0: api.Order.Items:type_name -> api.Item
	8,  // 1: api.CheckIfItemsInStockRequest.Items:type_name -> api.ItemsWithQuantity
	7,  // 2: api.CheckIfItemsInStockResponse.Items:type_name -> api.Item
	7,  // 3: api.GetItemsResponse.Items:type_name -> api.Item
	8,  // 4: api.CreateOrderRequest.Items:type_name -> api.ItemsWithQuantity
	9,  // 5: api.OrderService.CreateOrder:input_type -> api.CreateOrderRequest
	6,  // 6: api.OrderService.GetOrder:input_type -> api.GetOrderRequest
	0,  // 7: api.OrderService.UpdateOrder:input_type -> api.Order
	2,  // 8: api.StockService.CheckIfItemsInStock:input_type -> api.CheckIfItemsInStockRequest
	4,  // 9: api.StockService.GetItems:input_type -> api.GetItemsRequest
	0,  // 10: api.OrderService.CreateOrder:output_type -> api.Order
	0,  // 11: api.OrderService.GetOrder:output_type -> api.Order
	0,  // 12: api.OrderService.UpdateOrder:output_type -> api.Order
	3,  // 13: api.StockService.CheckIfItemsInStock:output_type -> api.CheckIfItemsInStockResponse
	5,  // 14: api.StockService.GetItems:output_type -> api.GetItemsResponse
	10, // [10:15] is the sub-list for method output_type
	5,  // [5:10] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_api_oms_proto_init() }
func file_api_oms_proto_init() {
	if File_api_oms_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_oms_proto_rawDesc), len(file_api_oms_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_api_oms_proto_goTypes,
		DependencyIndexes: file_api_oms_proto_depIdxs,
		MessageInfos:      file_api_oms_proto_msgTypes,
	}.Build()
	File_api_oms_proto = out.File
	file_api_oms_proto_goTypes = nil
	file_api_oms_proto_depIdxs = nil
}
