package product

import (
	"encoding/json"
	"github.com/go-ap/errors"
	"github.com/go-ap/filters"
	"strings"
	"time"

	"git.sr.ht/~mariusor/lw"
	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/fedbox/internal/cmd/ecommerce/common"
	"github.com/go-resty/resty/v2"
)

type ProductService struct {
	db      common.Storage
	baseURL string
	logger  lw.Logger
}

// NewProductService creates new ProductService instance
func NewProductService(db common.Storage, baseURL string, l lw.Logger) *ProductService {
	var target ProductService
	target.db = db
	target.baseURL = baseURL
	target.logger = l
	return &target
}

func (s *ProductService) GetProducts(caller vocab.Actor, token string) ([]Product, error) {
	//preparing the result product list
	var products = make([]Product, 0)

	//determining current actor ID
	actorIdParts := strings.Split(caller.ID.String(), "/")
	actorID := actorIdParts[len(actorIdParts)-1]

	//constructing url to get ActivityPub collection
	objectsUrl := s.baseURL + "/objects?type=Document&attributedTo=" + actorID
	s.logger.Infof("ProductService. Collection URL=%s", objectsUrl)

	//getting ActivityPub collection
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"").
		SetHeader("Authorization", token).
		Get(objectsUrl)
	if err != nil {
		s.logger.Errorf("product searching error", err)
		return products, err
	}

	//Analyzing the result
	if resp.IsSuccess() {
		//success
		s.logger.Infof("ProductService. Collection found. Result object:\n %s", resp.Body())

		//unmarshalling ActivityPub Collection json
		vocabCollection := vocab.OrderedCollectionNew(vocab.EmptyIRI)
		err = vocabCollection.UnmarshalJSON(resp.Body())
		if err != nil {
			s.logger.Errorf("parsing error", err)
			return products, err
		}
		s.logger.Infof("ProductService. Collection parsed successfully. Items count = %s", vocabCollection.TotalItems)

		//adding collection items to Product slice
		for _, item := range vocabCollection.OrderedItems {
			if !item.IsObject() {
				//skipping all not object items
				continue
			}
			var o *vocab.Object
			o, err = vocab.ToObject(item)
			if err != nil {
				//skipping an item in case of parsing error
				s.logger.Errorf("item to Product mapping error", err)
				continue
			}
			p := s.mapObjectToProduct(o)
			products = append(products, p)
		}
		s.logger.Infof("ProductService. Collection successfully mapped to Products. Products count = %s", len(products))

		return products, nil
	} else {
		//not success
		if resp.StatusCode() == 404 {
			//not found
			s.logger.Infof("ProductService. Product not found.")
			return nil, errors.NotFoundf("Products not found")
		} else {
			//some other problem...
			s.logger.Infof("Product not found. Status code=%s", resp.Status())
			return nil, errors.Newf("product not found")
		}
	}
}

func (s *ProductService) GetProduct(caller vocab.Actor, token string, productID string) (*Product, error) {
	//Constructing url to get ActivityPub object
	objectsUrl := s.baseURL + "/objects/" + productID + "?type=Document&attributedTo=" + caller.ID.String()
	s.logger.Infof("ProductService. Objects URL=%s", objectsUrl)

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"").
		SetHeader("Authorization", token).
		Get(objectsUrl)
	if err != nil {
		s.logger.Errorf("product searching error", err)
		return nil, err
	}

	//Analyzing the result
	if resp.IsSuccess() {
		//success
		s.logger.Infof("ProductService. Product found. Result product object:\n %s", resp.Body())

		//unmarshalling ActivityPub Object json
		vocabObject := vocab.ObjectNew(vocab.DocumentType)
		err = vocabObject.UnmarshalJSON(resp.Body())
		if err != nil {
			s.logger.Errorf("object parsing error", err)
			return nil, err
		}

		//mapping ActivityPub object to Product dto
		resultProduct := s.mapObjectToProduct(vocabObject)
		s.logger.Infof("ProductService. Result product object parsed successfully. Product.id=%s", resultProduct.Id)

		return &resultProduct, nil
	} else {
		//not success
		if resp.StatusCode() == 404 {
			//not found
			s.logger.Infof("ProductService. Product not found.")
			return nil, errors.NotFoundf("Product not found")
		} else {
			//some other problem...
			s.logger.Infof("Product not found. Status code=%s", resp.Status())
			return nil, errors.Newf("product not found")
		}
	}
}

// CreateProduct creates new Product object by posting Create Activity
func (s *ProductService) CreateProduct(caller vocab.Actor, token string, p Product) (vocab.Item, error) {
	//building ActivityPub object
	object, err := s.mapProductToObject(p, caller)
	if err != nil {
		s.logger.Errorf("product activity creation error", err)
		return nil, err
	}

	//wrapping Object to Create activity
	activity := vocab.CreateNew(vocab.EmptyIRI, object)
	activity.Actor = caller

	//marshaling activity to json
	activityJson, err := activity.MarshalJSON()
	if err != nil {
		s.logger.Errorf("create product activity marshaling error")
		return nil, err
	}

	//POST activity to current actor's outbox collection
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"").
		SetHeader("Authorization", token).
		SetBody(string(activityJson)).
		Post(caller.Outbox.GetLink().String())
	if err != nil {
		s.logger.Errorf("product creation error", err)
		return nil, err
	}
	s.logger.Debugf("ProductService. Product created. Result product object:\n %s", resp.Body())

	//unmarshalling ActivityPub Object json
	vocabObject := vocab.ObjectNew(vocab.ObjectType)
	err = vocabObject.UnmarshalJSON(resp.Body())
	if err != nil {
		s.logger.Errorf("object parsing error", err)
		return nil, err
	}

	return vocabObject, nil
}

// UpdateProduct updates the Product object by posting Update Activity
func (s *ProductService) UpdateProduct(caller vocab.Actor, token string, p Product) (vocab.Item, error) {
	//building ActivityPub object
	object, err := s.mapProductToObject(p, caller)
	if err != nil {
		s.logger.Errorf("product activity creation error", err)
		return nil, err
	}

	//wrapping Object to Create activity
	activity := vocab.UpdateNew(vocab.EmptyIRI, object)
	activity.Actor = caller

	//marshaling activity to json
	activityJson, err := activity.MarshalJSON()
	if err != nil {
		s.logger.Errorf("create product activity marshaling error")
		return nil, err
	}

	//POST activity to current actor's outbox collection
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"").
		SetHeader("Authorization", token).
		SetBody(string(activityJson)).
		Post(caller.Outbox.GetLink().String())
	if err != nil {
		s.logger.Errorf("product creation error", err)
		return nil, err
	}
	s.logger.Debugf("ProductService. Product created. Result product object:\n %s", resp.Body())

	//unmarshalling ActivityPub Object json
	vocabObject := vocab.ObjectNew(vocab.ObjectType)
	err = vocabObject.UnmarshalJSON(resp.Body())
	if err != nil {
		s.logger.Errorf("object parsing error", err)
		return nil, err
	}

	return vocabObject, nil
}

// DeleteProduct deletes the Product object by posting Delete Activity
func (s *ProductService) DeleteProduct(caller vocab.Actor, token string, productID string) error {
	//creating an Object with ID only
	object := vocab.ObjectNew(vocab.DocumentType)
	object.ID = vocab.IRI(s.baseURL + "/objects/" + productID)

	//wrapping Object to Delete activity
	activity := vocab.DeleteNew(vocab.EmptyIRI, object)
	activity.Actor = caller

	//marshaling activity to json
	activityJson, err := activity.MarshalJSON()
	if err != nil {
		s.logger.Errorf("create product activity marshaling error")
		return err
	}

	//POST activity to current actor's outbox collection
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"").
		SetHeader("Authorization", token).
		SetBody(string(activityJson)).
		Post(caller.Outbox.GetLink().String())
	if err != nil {
		s.logger.Errorf("product creation error", err)
		return err
	}
	s.logger.Debugf("ProductService. Product deleted. Result product object:\n %s", resp.Body())

	return nil
}

// parseResultProductObject parses json of ActivityPub Object type and puts it to Product dto
func (s *ProductService) mapObjectToProduct(o *vocab.Object) Product {
	// extracting product ID from ActivityPub Object ID
	productIdParts := strings.Split(o.ID.String(), "/")
	productID := productIdParts[len(productIdParts)-1]

	//mapping to dto
	p := Product{}
	p.Id = productID
	p.Name = o.Name.String()
	p.Summary = o.Summary.String()

	//collecting product images
	vocab.OnCollectionIntf(o.Attachment, func(col vocab.CollectionInterface) error {
		for _, it := range col.Collection() {
			vocab.OnObject(it, func(object *vocab.Object) error {
				productImage := Image{}
				productImage.Name = object.Name.String()
				productImage.Content = object.Content.String()
				if object.URL != nil {
					productImage.URL = object.URL.GetLink().String()
				}
				p.Images = append(p.Images, productImage)
				return nil
			})
		}
		return nil
	})

	//collecting product tags
	vocab.OnCollectionIntf(o.Tag, func(col vocab.CollectionInterface) error {
		for _, it := range col.Collection() {
			vocab.OnObject(it, func(object *vocab.Object) error {
				p.Tags = append(p.Tags, object.Name.First().String())
				return nil
			})
		}
		return nil
	})

	if o.Content.String() != "" {
		s.logger.Infof("ProductService. Product content = %s", o.Content.String()[1:len(o.Content.String())-1])

		//removing "" wrapper
		productContentString := o.Content.String()[1 : len(o.Content.String())-1]

		//unmarshal the content value to dto
		err := json.Unmarshal([]byte(productContentString), &p.Content)
		if err != nil {
			//just ignore the content in case of parsing errors
			s.logger.Warnf("parsing product content error")
		}
	} else {
		s.logger.Infof("ProductService. Product content is empty")
	}

	return p
}

// mapProductToObject builds ActivityPub object from given Product dto
func (s *ProductService) mapProductToObject(p Product, owner vocab.Actor) (*vocab.Object, error) {
	//creating Object
	o := vocab.ObjectNew(vocab.DocumentType)
	if p.Id != "" {
		o.ID = vocab.IRI(s.baseURL + "/objects/" + p.Id)
	}
	o.Published = time.Now()
	o.Name = vocab.DefaultNaturalLanguageValue(p.Name)
	o.Summary = vocab.DefaultNaturalLanguageValue(p.Summary)
	o.AttributedTo = owner.ID

	o.Attachment = s.prepareImages(p)

	tags := s.prepareTags(owner, p.Tags)
	if len(tags) > 0 {
		o.Tag = tags
	}
	//Object.Content property is a custom json string
	if len(p.Content) != 0 {
		s.logger.Infof("product content is not empty")
		//marshaling the content
		productContentData, err := json.Marshal(p.Content)
		if err != nil {
			s.logger.Errorf("parsing product content error")
			return nil, err
		}
		//store the content in 'Content' property wrapped in ""
		o.Content = vocab.DefaultNaturalLanguageValue("\"" + string(productContentData) + "\"")
		s.logger.Infof("product content = %s", o.Content.String())
	} else {
		s.logger.Infof("product content is empty")
	}

	return o, nil
}

func (s ProductService) prepareImages(p Product) vocab.ItemCollection {
	images := make(vocab.ItemCollection, 0)

	if len(p.Images) == 0 {
		s.logger.Infof("product contains bo images")
		return images
	}

	s.logger.Infof("product contains %d images", len(p.Images))
	for _, val := range p.Images {
		if len(val.URL) == 0 {
			continue
		}
		image := vocab.ObjectNew(vocab.ImageType)
		image.Name = vocab.DefaultNaturalLanguageValue(val.Name)
		image.Content = vocab.DefaultNaturalLanguageValue(val.Content)
		image.URL = vocab.IRI(val.URL)
		_ = images.Append(image)
	}

	return images
}

func (s *ProductService) prepareTags(owner vocab.Actor, src []string) vocab.ItemCollection {
	tags := make(vocab.ItemCollection, 0)

	existsTagFilter := filters.Filters{
		BaseURL:       vocab.IRI(s.baseURL),
		Authenticated: &owner,
		Type: filters.CompStrs{
			filters.CompStr{
				Str: string(vocab.ObjectType),
			},
		},
		IRI: vocab.IRI(s.baseURL + "/objects"),
		AttrTo: filters.CompStrs{
			filters.CompStr{
				Str: owner.ID.String(),
			},
		},
	}

	allObjects, _ := s.db.Load(existsTagFilter.GetLink())
	existsTags := map[string]*vocab.Object{}

	vocab.OnCollectionIntf(allObjects, func(col vocab.CollectionInterface) error {
		for _, it := range col.Collection() {
			vocab.OnObject(it, func(object *vocab.Object) error {
				existsTags[object.Name.First().Value.String()] = object
				return nil
			})
		}
		return nil
	})
	for _, t := range src {
		if el, ok := existsTags[t]; !ok {
			tag := vocab.ObjectNew("")
			tag.Name = vocab.NaturalLanguageValues{
				{Ref: vocab.NilLangRef, Value: vocab.Content(t)},
			}
			tag.AttributedTo = owner
			_ = tags.Append(tag)
		} else {
			_ = tags.Append(el)
		}
	}
	return tags
}
